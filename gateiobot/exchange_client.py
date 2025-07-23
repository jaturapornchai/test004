"""
Gate.io Exchange Client - ใช้ https://api.gateio.ws/api/v4
ตาม step.md specifications
"""

import os
import time
import hmac
import hashlib
import json
import requests
from typing import Dict, List, Optional, Any
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

class ExchangeClient:
    """Gate.io Exchange Client สำหรับ Futures Trading"""
    
    def __init__(self):
        """Initialize Exchange Client"""
        self.api_key = os.getenv('GATE_API_KEY')
        self.api_secret = os.getenv('GATE_API_SECRET')
        
        if not self.api_key or not self.api_secret:
            raise ValueError("❌ ไม่พบ Gate.io API Key หรือ Secret ใน .env file")
        
        # API endpoints ตาม step.md
        self.base_url = "https://api.gateio.ws/api/v4"
        self.futures_base = "/futures/usdt"
        
        # Session for connection reuse
        self.session = requests.Session()
        
        print("✅ เชื่อมต่อ Gate.io Exchange Client สำเร็จ")
    
    def _generate_signature(self, method: str, url: str, query_string: str = "", payload_string: str = "") -> Dict[str, str]:
        """สร้าง signature สำหรับ authentication ตาม Gate.io API v4 specs"""
        try:
            timestamp = str(int(time.time()))
            
            # สร้าง payload hash (ใช้ empty string ถ้าไม่มี payload)
            if payload_string:
                payload_hash = hashlib.sha512(payload_string.encode('utf-8')).hexdigest()
            else:
                payload_hash = hashlib.sha512(b'').hexdigest()
            
            # สร้าง message สำหรับ signing ตาม Gate.io format
            sign_string = f"{method}\n{url}\n{query_string}\n{payload_hash}\n{timestamp}"
            
            # Generate signature
            signature = hmac.new(
                self.api_secret.encode('utf-8'),
                sign_string.encode('utf-8'),
                hashlib.sha512
            ).hexdigest()
            
            return {
                "KEY": self.api_key,
                "Timestamp": timestamp,
                "SIGN": signature
            }
        except Exception as e:
            print("❌ เกิดข้อผิดพลาดในการสร้าง signature: " + str(e))
            return {}
    
    def _make_request(self, method: str, endpoint: str, params: Optional[Dict] = None, data: Optional[Dict] = None) -> Optional[Any]:
        """ส่ง HTTP request ไปยัง Gate.io API ตาม v4 specification"""
        try:
            # สร้าง full URL - ถ้า endpoint ไม่ขึ้นต้นด้วย /api/v4 ให้เพิ่ม
            if endpoint.startswith('/api/v4'):
                full_url = self.base_url + endpoint
            else:
                full_url = self.base_url + "/api/v4" + endpoint
            
            query_string = ""
            body = ""
            
            # สร้าง query string สำหรับ GET parameters
            if params and method == "GET":
                query_string = "&".join([f"{k}={v}" for k, v in params.items()])
            
            # สร้าง body สำหรับ POST data
            if data and method == "POST":
                body = json.dumps(data)
            
            # Generate headers ตาม Gate.io API v4 specification
            headers = self._generate_signature(method, endpoint, query_string, body)
            headers["Content-Type"] = "application/json"
            headers["Accept"] = "application/json"
            
            # Make request
            if method == "GET":
                if query_string:
                    full_url += "?" + query_string
                response = self.session.get(full_url, headers=headers, timeout=10)
            elif method == "POST":
                response = self.session.post(full_url, headers=headers, data=body, timeout=10)
            elif method == "DELETE":
                response = self.session.delete(full_url, headers=headers, timeout=10)
            else:
                print("❌ HTTP method ไม่รองรับ: " + method)
                return None
            
            if response.status_code == 200 or response.status_code == 201:
                return response.json()
            else:
                print("❌ API Error " + str(response.status_code) + ": " + response.text)
                return None
                
        except Exception as e:
            print("❌ เกิดข้อผิดพลาดในการเรียก API: " + str(e))
            return None
    
    def test_connection(self) -> bool:
        """ทดสอบการเชื่อมต่อ API"""
        try:
            result = self._make_request("GET", "/api/v4/futures/usdt/contracts")
            if result:
                print("✅ ทดสอบการเชื่อมต่อ Gate.io API สำเร็จ")
                return True
            else:
                print("❌ ไม่สามารถเชื่อมต่อ Gate.io API ได้")
                return False
        except Exception as e:
            print("❌ เกิดข้อผิดพลาดในการทดสอบการเชื่อมต่อ: " + str(e))
            return False
    
    def get_exchange(self):
        """Return self สำหรับ compatibility"""
        return self
    
    # ============= Market Data Methods =============
    
    def load_markets(self) -> Dict[str, Any]:
        """ดึงข้อมูล markets (futures contracts)"""
        try:
            contracts = self._make_request("GET", "/futures/usdt/contracts")
            if contracts:
                markets = {}
                for contract in contracts:
                    symbol = contract.get('name', '').replace('_', '/') + ':USDT'
                    markets[symbol] = {
                        'id': contract.get('name'),
                        'symbol': symbol,
                        'base': contract.get('underlying', '').replace('_USDT', ''),
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
                        'contractSize': float(contract.get('quanto_multiplier', 1)),
                        'limits': {
                            'amount': {
                                'min': float(contract.get('order_size_min', 1)),
                                'max': float(contract.get('order_size_max', 1000000))
                            }
                        },
                        'precision': {
                            'amount': 8,
                            'price': 8
                        },
                        'info': contract
                    }
                print("✅ โหลด " + str(len(markets)) + " markets สำเร็จ")
                return markets
            return {}
        except Exception as e:
            print("❌ เกิดข้อผิดพลาดในการดึงข้อมูล markets: " + str(e))
            return {}
    
    @property  
    def markets(self) -> Dict[str, Any]:
        """Property สำหรับ markets (lazy loading)"""
        if not hasattr(self, '_markets'):
            self._markets = self.load_markets()
        return self._markets
    
    def fetch_ohlcv(self, symbol: str, timeframe: str = '1h', limit: int = 120) -> List[List]:
        """ดึงข้อมูล OHLCV"""
        try:
            # แปลง symbol format
            contract = symbol.replace('/USDT:USDT', '_USDT').replace('/', '_')
            
            # แปลง timeframe
            interval_map = {'1h': '1h', '4h': '4h', '1d': '1d'}
            interval = interval_map.get(timeframe, '1h')
            
            params = {
                'contract': contract,
                'interval': interval,
                'limit': limit
            }
            
            candles = self._make_request("GET", "/futures/usdt/candlesticks", params=params)
            
            if candles:
                ohlcv = []
                for candle in candles:
                    ohlcv.append([
                        int(candle['t']) * 1000,  # timestamp
                        float(candle['o']),       # open
                        float(candle['h']),       # high
                        float(candle['l']),       # low
                        float(candle['c']),       # close
                        float(candle['v'])        # volume
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
            
            tickers = self._make_request("GET", "/futures/usdt/tickers", params={'contract': contract})
            
            if tickers and len(tickers) > 0:
                ticker = tickers[0]
                return {
                    'symbol': symbol,
                    'last': float(ticker.get('last', 0)),
                    'bid': float(ticker.get('highest_bid', 0)),
                    'ask': float(ticker.get('lowest_ask', 0)),
                    'high': float(ticker.get('high_24h', 0)),
                    'low': float(ticker.get('low_24h', 0)),
                    'quoteVolume': float(ticker.get('quote_volume_24h', 0)),
                    'baseVolume': float(ticker.get('base_volume_24h', 0)),
                    'change': float(ticker.get('change_percentage', 0)),
                    'info': ticker
                }
            return None
        except Exception as e:
            print("❌ เกิดข้อผิดพลาดในการดึงข้อมูล ticker " + symbol + ": " + str(e))
            return None
    
    # ============= Account Methods =============
    
    def fetch_balance(self) -> Dict[str, Any]:
        """ดึงข้อมูล balance จาก Futures account"""
        try:
            # ใช้ futures balance endpoint
            result = self._make_request("GET", "/futures/usdt/accounts")
            
            if result:
                return {
                    'USDT': {
                        'free': float(result.get('available', 0)),
                        'used': float(result.get('position_margin', 0)),
                        'total': float(result.get('total', 0))
                    },
                    'info': result
                }
            return {}
        except Exception as e:
            print("❌ เกิดข้อผิดพลาดในการดึงข้อมูล balance: " + str(e))
            return {}
    
    def fetch_positions(self) -> List[Dict]:
        """ดึงข้อมูล positions"""
        try:
            positions = self._make_request("GET", "/futures/usdt/positions")
            
            if positions:
                result = []
                for pos in positions:
                    if float(pos.get('size', 0)) != 0:  # เฉพาะ position ที่เปิดอยู่
                        symbol = pos.get('contract', '').replace('_', '/') + ':USDT'
                        result.append({
                            'symbol': symbol,
                            'contracts': float(pos.get('size', 0)),
                            'contractSize': 1,
                            'unrealizedPnl': float(pos.get('unrealised_pnl', 0)),
                            'percentage': float(pos.get('unrealised_pnl', 0)) / float(pos.get('margin', 1)) * 100 if float(pos.get('margin', 0)) > 0 else 0,
                            'entryPrice': float(pos.get('entry_price', 0)),
                            'notional': float(pos.get('value', 0)),
                            'timestamp': int(time.time() * 1000),
                            'side': 'long' if float(pos.get('size', 0)) > 0 else 'short',
                            'info': pos
                        })
                return result
            return []
        except Exception as e:
            print("❌ เกิดข้อผิดพลาดในการดึงข้อมูล positions: " + str(e))
            return []
    
    def create_order(self, symbol: str, type: str, side: str, amount: float, 
                    price: Optional[float] = None, reduce_only: bool = False, **params) -> Dict[str, Any]:
        """
        สร้าง order บน Gate.io Futures ตาม API v4 specification
        
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
            
            # สร้าง order data ตามตัวอย่าง
            order_data = {
                'contract': contract,
                'size': size,
                'price': str(price) if price else "0",  # "0" สำหรับ market order
                'tif': 'ioc'  # Immediate or Cancel สำหรับ market order
            }
            
            # เพิ่ม reduce_only ถ้าต้องการ
            if reduce_only:
                order_data['reduce_only'] = True
            
            # ส่ง request ตาม API v4 specification
            endpoint = "/api/v4/futures/usdt/orders"
            result = self._make_request('POST', endpoint, data=order_data)
            
            if result and not result.get('error'):
                return {
                    'id': result.get('id'),
                    'symbol': symbol,
                    'type': type,
                    'side': side,
                    'amount': amount,
                    'price': price,
                    'status': 'open',
                    'info': result
                }
            else:
                error_msg = result.get('message', 'Unknown error') if result else 'No response from API'
                return {'error': error_msg}
                
        except Exception as e:
            error_msg = str(e)
            print(f"❌ เกิดข้อผิดพลาดในการสร้าง order: {error_msg}")
            return {'error': error_msg}
