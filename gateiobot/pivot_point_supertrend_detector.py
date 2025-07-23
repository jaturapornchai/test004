"""
Pivot Point SuperTrend Detector ตาม step.md specifications
โดยใช้ Pivot Point SuperTrend + EMA100 เท่านั้น (ไม่ใช้ RSI ตาม specifications)
"""

import numpy as np
import pandas as pd
from typing import Dict, List, Tuple, Optional
import logging

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class PivotPointSuperTrendDetector:
    """
    Pivot Point SuperTrend Detector ตาม step.md specifications
    - ใช้ Pivot Point SuperTrend + EMA100 เท่านั้น
    - ไม่ใช้ RSI (removed ตาม specifications)
    - AUTO OPEN Mode: SuperTrend Confidence ≥ 75%
    - AI Mode: SuperTrend Confidence < 75%
    """
    
    def __init__(self, 
                 pivot_length: int = 7, 
                 multiplier: float = 3.0,
                 ema_length: int = 100):
        """
        กำหนดค่า Pivot Point SuperTrend
        
        Args:
            pivot_length: ความยาว pivot point (default: 7)
            multiplier: ตัวคูณ SuperTrend (default: 3.0)
            ema_length: ความยาว EMA (default: 100)
        """
        self.pivot_length = pivot_length
        self.multiplier = multiplier
        self.ema_length = ema_length
        
        logger.info(f"🎯 เริ่มต้น Pivot Point SuperTrend Detector")
        logger.info(f"   Pivot Length: {pivot_length}")
        logger.info(f"   Multiplier: {multiplier}")
        logger.info(f"   EMA Length: {ema_length}")
    
    def calculate_pivot_points(self, high: np.array, low: np.array) -> Tuple[np.array, np.array]:
        """
        คำนวณ Pivot Points
        
        Args:
            high: ราคาสูงสุด
            low: ราคาต่ำสุด
            
        Returns:
            Tuple[pivot_high, pivot_low]
        """
        pivot_high = np.full(len(high), np.nan)
        pivot_low = np.full(len(low), np.nan)
        
        for i in range(self.pivot_length, len(high) - self.pivot_length):
            # Pivot High
            is_pivot_high = True
            for j in range(self.pivot_length):
                if high[i] <= high[i-j-1] or high[i] <= high[i+j+1]:
                    is_pivot_high = False
                    break
            
            if is_pivot_high:
                pivot_high[i] = high[i]
            
            # Pivot Low
            is_pivot_low = True
            for j in range(self.pivot_length):
                if low[i] >= low[i-j-1] or low[i] >= low[i+j+1]:
                    is_pivot_low = False
                    break
            
            if is_pivot_low:
                pivot_low[i] = low[i]
        
        return pivot_high, pivot_low
    
    def calculate_supertrend(self, high: np.array, low: np.array, close: np.array) -> Tuple[np.array, np.array]:
        """
        คำนวณ Pivot Point SuperTrend
        
        Args:
            high: ราคาสูงสุด
            low: ราคาต่ำสุด
            close: ราคาปิด
            
        Returns:
            Tuple[supertrend, trend_direction]
        """
        # คำนวณ pivot points
        pivot_high, pivot_low = self.calculate_pivot_points(high, low)
        
        # คำนวณ ATR จาก pivot points
        tr1 = high - low
        tr2 = np.abs(high - np.roll(close, 1))
        tr3 = np.abs(low - np.roll(close, 1))
        true_range = np.maximum(tr1, np.maximum(tr2, tr3))
        
        # ใช้ pivot length สำหรับ ATR
        atr = np.full(len(close), np.nan)
        for i in range(self.pivot_length, len(close)):
            atr[i] = np.mean(true_range[i-self.pivot_length+1:i+1])
        
        # คำนวณ basic upper และ lower bands
        hl2 = (high + low) / 2
        basic_upper = hl2 + (self.multiplier * atr)
        basic_lower = hl2 - (self.multiplier * atr)
        
        # คำนวณ final upper และ lower bands
        final_upper = np.full(len(close), np.nan)
        final_lower = np.full(len(close), np.nan)
        
        for i in range(1, len(close)):
            if not np.isnan(basic_upper[i]):
                final_upper[i] = basic_upper[i] if np.isnan(final_upper[i-1]) or basic_upper[i] < final_upper[i-1] or close[i-1] > final_upper[i-1] else final_upper[i-1]
            
            if not np.isnan(basic_lower[i]):
                final_lower[i] = basic_lower[i] if np.isnan(final_lower[i-1]) or basic_lower[i] > final_lower[i-1] or close[i-1] < final_lower[i-1] else final_lower[i-1]
        
        # คำนวณ SuperTrend
        supertrend = np.full(len(close), np.nan)
        trend_direction = np.full(len(close), np.nan)  # 1 = bullish, -1 = bearish
        
        for i in range(1, len(close)):
            if not np.isnan(final_upper[i]) and not np.isnan(final_lower[i]):
                if np.isnan(supertrend[i-1]):
                    supertrend[i] = final_upper[i] if close[i] <= final_upper[i] else final_lower[i]
                    trend_direction[i] = -1 if close[i] <= final_upper[i] else 1
                else:
                    if supertrend[i-1] == final_upper[i-1] and close[i] > final_upper[i]:
                        supertrend[i] = final_lower[i]
                        trend_direction[i] = 1
                    elif supertrend[i-1] == final_lower[i-1] and close[i] < final_lower[i]:
                        supertrend[i] = final_upper[i]
                        trend_direction[i] = -1
                    else:
                        supertrend[i] = supertrend[i-1]
                        trend_direction[i] = trend_direction[i-1]
        
        return supertrend, trend_direction
    
    def calculate_ema(self, data: np.array, length: int) -> np.array:
        """
        คำนวณ EMA
        
        Args:
            data: ข้อมูลราคา
            length: ความยาว EMA
            
        Returns:
            EMA values
        """
        ema = np.full(len(data), np.nan)
        alpha = 2.0 / (length + 1)
        
        # หาจุดเริ่มต้นที่ไม่เป็น NaN
        start_idx = 0
        for i in range(len(data)):
            if not np.isnan(data[i]):
                ema[i] = data[i]
                start_idx = i
                break
        
        # คำนวณ EMA
        for i in range(start_idx + 1, len(data)):
            if not np.isnan(data[i]):
                ema[i] = alpha * data[i] + (1 - alpha) * ema[i-1]
        
        return ema
    
    def calculate_confidence(self, close: np.array, supertrend: np.array, 
                           trend_direction: np.array, ema100: np.array) -> float:
        """
        คำนวณ SuperTrend Confidence
        
        Args:
            close: ราคาปิด
            supertrend: SuperTrend values
            trend_direction: ทิศทาง trend
            ema100: EMA100 values
            
        Returns:
            Confidence percentage (0-100)
        """
        if len(close) < 10:
            return 0.0
        
        confidence = 0.0
        
        # ตรวจสอบ trend consistency (40%)
        recent_trends = trend_direction[-10:]  # 10 แท่งล่าสุด
        if not np.isnan(recent_trends).all():
            trend_consistency = np.sum(recent_trends == recent_trends[-1]) / len(recent_trends)
            confidence += trend_consistency * 40
        
        # ตรวจสอบ distance from SuperTrend (30%)
        current_price = close[-1]
        current_supertrend = supertrend[-1]
        if not np.isnan(current_supertrend):
            distance_ratio = abs(current_price - current_supertrend) / current_price
            distance_score = min(distance_ratio * 10, 1.0)  # ยิ่งไกลยิ่งดี (แต่ไม่เกิน 100%)
            confidence += distance_score * 30
        
        # ตรวจสอบ EMA100 confirmation (30%)
        current_ema100 = ema100[-1]
        if not np.isnan(current_ema100):
            if trend_direction[-1] == 1 and current_price > current_ema100:  # Bullish
                confidence += 30
            elif trend_direction[-1] == -1 and current_price < current_ema100:  # Bearish
                confidence += 30
        
        return min(confidence, 100.0)
    
    def analyze(self, ohlcv_data: List[List]) -> Dict:
        """
        วิเคราะห์ Pivot Point SuperTrend
        
        Args:
            ohlcv_data: ข้อมูล OHLCV [[timestamp, open, high, low, close, volume], ...]
            
        Returns:
            ผลการวิเคราะห์
        """
        try:
            if len(ohlcv_data) < max(self.pivot_length * 2, self.ema_length):
                logger.warning(f"ข้อมูลไม่เพียงพอ: {len(ohlcv_data)} แท่ง (ต้องการ {max(self.pivot_length * 2, self.ema_length)} แท่ง)")
                return {
                    'signal': 'HOLD',
                    'confidence': 0.0,
                    'trend_direction': 0,
                    'mode': 'INSUFFICIENT_DATA',
                    'current_price': 0.0,
                    'supertrend_value': 0.0,
                    'ema100_value': 0.0,
                    'analysis': 'ข้อมูลไม่เพียงพอสำหรับการวิเคราะห์'
                }
            
            # แปลงข้อมูลเป็น numpy arrays
            df = pd.DataFrame(ohlcv_data, columns=['timestamp', 'open', 'high', 'low', 'close', 'volume'])
            
            high = df['high'].values.astype(float)
            low = df['low'].values.astype(float)
            close = df['close'].values.astype(float)
            
            # คำนวณ SuperTrend
            supertrend, trend_direction = self.calculate_supertrend(high, low, close)
            
            # คำนวณ EMA100
            ema100 = self.calculate_ema(close, self.ema_length)
            
            # คำนวณ confidence
            confidence = self.calculate_confidence(close, supertrend, trend_direction, ema100)
            
            # กำหนด signal
            current_trend = trend_direction[-1] if not np.isnan(trend_direction[-1]) else 0
            current_price = close[-1]
            current_supertrend = supertrend[-1] if not np.isnan(supertrend[-1]) else 0
            current_ema100 = ema100[-1] if not np.isnan(ema100[-1]) else 0
            
            # กำหนด signal ตาม trend
            if current_trend == 1:
                signal = 'BUY'
            elif current_trend == -1:
                signal = 'SELL'
            else:
                signal = 'HOLD'
            
            # กำหนด mode ตาม confidence
            mode = 'AUTO_OPEN' if confidence >= 75 else 'AI'
            
            # สร้าง analysis text
            analysis = f"SuperTrend: {'Bullish' if current_trend == 1 else 'Bearish' if current_trend == -1 else 'Neutral'}, "
            analysis += f"Confidence: {confidence:.1f}%, "
            analysis += f"Mode: {mode}"
            
            result = {
                'signal': signal,
                'confidence': confidence,
                'trend_direction': int(current_trend),
                'mode': mode,
                'current_price': float(current_price),
                'supertrend_value': float(current_supertrend),
                'ema100_value': float(current_ema100),
                'analysis': analysis
            }
            
            # ปิด logging เพื่อป้องกัน duplicate output (app.py จะ handle การ print)
            # logger.info(f"📊 SuperTrend Analysis: {signal} | Confidence: {confidence:.1f}% | Mode: {mode}")
            
            return result
            
        except Exception as e:
            logger.error(f"❌ เกิดข้อผิดพลาดในการวิเคราะห์: {str(e)}")
            return {
                'signal': 'HOLD',
                'confidence': 0.0,
                'trend_direction': 0,
                'mode': 'ERROR',
                'current_price': 0.0,
                'supertrend_value': 0.0,
                'ema100_value': 0.0,
                'analysis': f'เกิดข้อผิดพลาด: {str(e)}'
            }
