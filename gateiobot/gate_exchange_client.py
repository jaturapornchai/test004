"""
Gate.io Exchange Client - ใช้ Official Gate.io Python Library
เขียนใหม่ทั้งหมดตาม step.md specifications
"""

import os
import gate_api
from gate_api.exceptions import ApiException, GateApiException
from typing import Dict, List, Optional, Any
import time

class GateExchangeClient:
    """Gate.io Exchange Client ใช้ Official Library"""
    
    def __init__(self):
        """เริ่มต้น Gate.io client ด้วย official library"""
        try:
            # อ่าน API credentials จาก .env
            from dotenv import load_dotenv
            load_dotenv()
            
            self.api_key = os.getenv('GATE_API_KEY')
            self.api_secret = os.getenv('GATE_API_SECRET')
            
            if not self.api_key or not self.api_secret:
                raise ValueError("❌ ไม่พบ GATE_API_KEY หรือ GATE_API_SECRET ใน .env file")
            
            # สร้าง configuration
            self.configuration = gate_api.Configuration(
                host="https://api.gateio.ws/api/v4",
                key=self.api_key,
                secret=self.api_secret
            )
            
            # สร้าง API client
            self.api_client = gate_api.ApiClient(self.configuration)
            self.futures_api = gate_api.FuturesApi(self.api_client)
            
            print("✅ เชื่อมต่อ Gate.io Exchange Client สำเร็จ")
            
        except Exception as e:
            print("❌ เกิดข้อผิดพลาดในการเชื่อมต่อ Gate.io: " + str(e))
            raise
    
    def test_connection(self) -> bool:
        """ทดสอบการเชื่อมต่อ API"""
        try:
            # ทดสอบด้วยการดึง contracts
            contracts = self.futures_api.list_futures_contracts('usdt')
            if contracts and len(contracts) > 0:
                print("✅ ทดสอบการเชื่อมต่อ Gate.io API สำเร็จ")
                return True
            else:
                print("❌ ไม่สามารถเชื่อมต่อ Gate.io API ได้")
                return False
        except Exception as e:
            print("❌ เกิดข้อผิดพลาดในการทดสอบการเชื่อมต่อ: " + str(e))
            return False
    
    def load_markets(self) -> Dict[str, Dict]:
        """โหลดข้อมูล markets"""
        try:
            contracts = self.futures_api.list_futures_contracts('usdt')
            markets = {}
            
            for contract in contracts:
                if contract.name.endswith('_USDT'):
                    symbol = contract.name.replace('_', '/') + ':USDT'
                    markets[symbol] = {
                        'id': contract.name,
                        'symbol': symbol,
                        'base': contract.name.replace('_USDT', ''),
                        'quote': 'USDT',
                        'settle': 'USDT',
                        'type': 'swap',
                        'spot': False,
                        'margin': False,
                        'swap': True,
                        'future': False,
                        'option': False,
                        'active': True,
                        'contract': True,
                        'linear': True,
                        'inverse': False,
                        'contractSize': float(contract.quanto_multiplier) if contract.quanto_multiplier else 1.0,
                        'precision': {
                            'amount': 8,
                            'price': 8
                        },
                        'limits': {
                            'amount': {
                                'min': float(contract.order_size_min) if contract.order_size_min else 1e-8,
                                'max': float(contract.order_size_max) if contract.order_size_max else None
                            }
                        },
                        'info': contract.to_dict()
                    }
            
            print("✅ โหลด " + str(len(markets)) + " markets สำเร็จ")
            return markets
            
        except Exception as e:
            print("❌ เกิดข้อผิดพลาดในการโหลด markets: " + str(e))
            return {}
    
    def fetch_balance(self) -> Dict[str, Any]:
        """ดึงข้อมูล balance"""
        try:
            account = self.futures_api.list_futures_accounts('usdt')
            
            if account and len(account) > 0:
                acc = account[0]
                return {
                    'USDT': {
                        'free': float(acc.available) if acc.available else 0.0,
                        'used': float(acc.position_margin) if acc.position_margin else 0.0,
                        'total': float(acc.total) if acc.total else 0.0
                    },
                    'info': acc.to_dict()
                }
            return {}
            
        except Exception as e:
            print("❌ เกิดข้อผิดพลาดในการดึงข้อมูล balance: " + str(e))
            return {}
    
    def fetch_positions(self) -> List[Dict]:
        """ดึงข้อมูล positions"""
        try:
            positions = self.futures_api.list_positions('usdt')
            result = []
            
            for pos in positions:
                if float(pos.size) != 0:  # เฉพาะ position ที่เปิดอยู่
                    symbol = pos.contract.replace('_', '/') + ':USDT'
                    result.append({
                        'symbol': symbol,
                        'contracts': float(pos.size),
                        'contractSize': 1.0,
                        'unrealizedPnl': float(pos.unrealised_pnl) if pos.unrealised_pnl else 0.0,
                        'percentage': float(pos.unrealised_pnl) / float(pos.margin) * 100 if pos.margin and float(pos.margin) > 0 else 0.0,
                        'entryPrice': float(pos.entry_price) if pos.entry_price else 0.0,
                        'notional': float(pos.value) if pos.value else 0.0,
                        'timestamp': int(time.time() * 1000),
                        'side': 'long' if float(pos.size) > 0 else 'short',
                        'info': pos.to_dict()
                    })
            
            return result
            
        except Exception as e:
            print("❌ เกิดข้อผิดพลาดในการดึงข้อมูล positions: " + str(e))
            return []
    
    def fetch_ohlcv(self, symbol: str, timeframe: str = '1h', limit: int = 120) -> List[List]:
        """ดึงข้อมูล OHLCV"""
        try:
            # แปลง symbol format
            contract = symbol.replace('/USDT:USDT', '_USDT').replace('/', '_')
            
            # แปลง timeframe
            interval_map = {'1h': '1h', '4h': '4h', '1d': '1d'}
            interval = interval_map.get(timeframe, '1h')
            
            # ดึงข้อมูล candlesticks
            candles = self.futures_api.list_futures_candlesticks(
                settle='usdt',
                contract=contract,
                _from=None,
                to=None,
                limit=limit,
                interval=interval
            )
            
            if candles:
                ohlcv = []
                for candle in candles:
                    ohlcv.append([
                        int(candle.t) * 1000,  # timestamp
                        float(candle.o),       # open
                        float(candle.h),       # high
                        float(candle.l),       # low
                        float(candle.c),       # close
                        float(candle.v)        # volume
                    ])
                return sorted(ohlcv, key=lambda x: x[0])  # เรียงตาม timestamp
            return []
            
        except Exception as e:
            print("❌ เกิดข้อผิดพลาดในการดึงข้อมูล OHLCV " + symbol + ": " + str(e))
            return []
    
    def fetch_ticker(self, symbol: str) -> Optional[Dict]:
        """ดึงข้อมูล ticker"""
        try:
            contract = symbol.replace('/USDT:USDT', '_USDT').replace('/', '_')
            
            tickers = self.futures_api.list_futures_tickers('usdt', contract=contract)
            
            if tickers and len(tickers) > 0:
                ticker = tickers[0]
                return {
                    'symbol': symbol,
                    'last': float(ticker.last),
                    'bid': float(ticker.highest_bid) if ticker.highest_bid else 0.0,
                    'ask': float(ticker.lowest_ask) if ticker.lowest_ask else 0.0,
                    'high': float(ticker.high_24h) if ticker.high_24h else 0.0,
                    'low': float(ticker.low_24h) if ticker.low_24h else 0.0,
                    'quoteVolume': float(ticker.quote_volume_24h) if ticker.quote_volume_24h else 0.0,
                    'baseVolume': float(ticker.base_volume_24h) if ticker.base_volume_24h else 0.0,
                    'change': float(ticker.change_percentage) if ticker.change_percentage else 0.0,
                    'info': ticker.to_dict()
                }
            return None
            
        except Exception as e:
            print("❌ เกิดข้อผิดพลาดในการดึงข้อมูล ticker " + symbol + ": " + str(e))
            return None
    
    def create_order(self, symbol: str, type: str, side: str, amount: float, 
                    price: Optional[float] = None, reduce_only: bool = False, **params) -> Dict[str, Any]:
        """
        สร้าง order บน Gate.io Futures ด้วย Official API
        
        Args:
            symbol: สัญลักษณ์เหรียญ เช่น BTC/USDT:USDT
            type: ประเภท order ('market', 'limit')
            side: ทิศทาง ('buy', 'sell')
            amount: จำนวนเหรียญ
            price: ราคา (สำหรับ limit order)
            reduce_only: ลดขนาด position เท่านั้น
            
        Returns:
            Dict: ผลลัพธ์การสร้าง order
        """
        try:
            # แปลง symbol จาก BTC/USDT:USDT เป็น BTC_USDT สำหรับ Gate.io
            contract = symbol.replace('/USDT:USDT', '_USDT').replace('/', '_')
            
            # คำนวณ size ตาม Gate.io format: + สำหรับ long, - สำหรับ short
            size = int(amount) if side == 'buy' else int(-amount)
            
            # สร้าง FuturesOrder object
            order = gate_api.FuturesOrder(
                contract=contract,
                size=size,
                price=str(price) if price else "0",  # "0" สำหรับ market order
                tif='ioc' if type == 'market' else 'gtc',  # IOC สำหรับ market, GTC สำหรับ limit
                reduce_only=reduce_only
            )
            
            # ส่ง order ผ่าน official API
            result = self.futures_api.create_futures_order(settle='usdt', futures_order=order)
            
            if result:
                return {
                    'id': str(result.id),
                    'symbol': symbol,
                    'type': type,
                    'side': side,
                    'amount': amount,
                    'price': price,
                    'status': 'open',
                    'info': result.to_dict()
                }
            else:
                return {'error': 'No response from API'}
                
        except GateApiException as ex:
            error_msg = "Gate API Error - Label: " + str(ex.label) + ", Message: " + str(ex.message)
            print("❌ " + error_msg)
            return {'error': error_msg}
        except ApiException as e:
            error_msg = "API Exception: " + str(e)
            print("❌ " + error_msg)
            return {'error': error_msg}
        except Exception as e:
            error_msg = "Unexpected error: " + str(e)
            print("❌ เกิดข้อผิดพลาดในการสร้าง order: " + error_msg)
            return {'error': error_msg}
    
    def close_position(self, symbol: str, side: Optional[str] = None) -> Dict[str, Any]:
        """ปิด position"""
        try:
            contract = symbol.replace('/USDT:USDT', '_USDT').replace('/', '_')
            
            # สร้าง order ปิด position
            order = gate_api.FuturesOrder(
                contract=contract,
                size=0,  # size = 0 สำหรับปิด position
                price="0",  # market order
                tif='ioc',
                close=True,  # close position
                reduce_only=True
            )
            
            result = self.futures_api.create_futures_order(settle='usdt', futures_order=order)
            
            if result:
                return {
                    'id': str(result.id),
                    'symbol': symbol,
                    'status': 'closed',
                    'info': result.to_dict()
                }
            else:
                return {'error': 'No response from API'}
                
        except Exception as e:
            error_msg = str(e)
            print("❌ เกิดข้อผิดพลาดในการปิด position: " + error_msg)
            return {'error': error_msg}
    
    def set_leverage(self, symbol: str, leverage: int = 5) -> bool:
        """ตั้งค่า leverage สำหรับ contract"""
        try:
            contract = symbol.replace('/USDT:USDT', '_USDT').replace('/', '_')
            
            # ตั้งค่า leverage
            self.futures_api.update_position_leverage(
                settle='usdt',
                contract=contract,
                leverage=str(leverage)
            )
            
            print(f"✅ ตั้ง leverage {leverage}x สำหรับ {symbol} สำเร็จ")
            return True
            
        except Exception as e:
            print(f"❌ เกิดข้อผิดพลาดในการตั้ง leverage {symbol}: {str(e)}")
            return False
    
    def set_margin_mode(self, symbol: str, margin_mode: str = 'isolated') -> bool:
        """ตั้งค่า margin mode สำหรับ contract"""
        try:
            contract = symbol.replace('/USDT:USDT', '_USDT').replace('/', '_')
            
            # ตั้งค่า margin mode (isolated = True, cross = False)
            is_isolated = margin_mode.lower() == 'isolated'
            
            self.futures_api.update_position_margin(
                settle='usdt',
                contract=contract,
                change=0,  # ไม่เพิ่มหรือลด margin แค่เปลี่ยน mode
                type='isolated' if is_isolated else 'cross'
            )
            
            print(f"✅ ตั้ง margin mode {margin_mode} สำหรับ {symbol} สำเร็จ")
            return True
            
        except Exception as e:
            print(f"❌ เกิดข้อผิดพลาดในการตั้ง margin mode {symbol}: {str(e)}")
            return False
    
    def ensure_position_settings(self, symbol: str, leverage: int = 5, margin_mode: str = 'isolated') -> bool:
        """ตรวจสอบและตั้งค่า position settings ให้ถูกต้องก่อนเปิด position"""
        try:
            print(f"🔧 ตรวจสอบการตั้งค่า {symbol}...")
            
            # ตั้งค่า leverage
            leverage_success = self.set_leverage(symbol, leverage)
            
            # ตั้งค่า margin mode  
            margin_success = self.set_margin_mode(symbol, margin_mode)
            
            if leverage_success and margin_success:
                print(f"✅ การตั้งค่า {symbol} สมบูรณ์ (Leverage: {leverage}x, Margin: {margin_mode})")
                return True
            else:
                print(f"⚠️ การตั้งค่า {symbol} มีปัญหาบางส่วน")
                return False
                
        except Exception as e:
            print(f"❌ เกิดข้อผิดพลาดในการตั้งค่า {symbol}: {str(e)}")
            return False
    
    def get_exchange(self):
        """Return self สำหรับ compatibility"""
        return self

if __name__ == "__main__":
    # ทดสอบการเชื่อมต่อ
    try:
        client = GateExchangeClient()
        if client.test_connection():
            print("✅ การทดสอบสำเร็จ")
        else:
            print("❌ การทดสอบล้มเหลว")
    except Exception as e:
        print("❌ เกิดข้อผิดพลาด: " + str(e))
