# system-broadcast-agent

# 📡 system-broadcast-agent (with grandcat/zeroconf)

로컬 네트워크 내에서 UDP 기반 mDNS (Bonjour) 프로토콜을 사용하여 데이터를 **송수신**하고, **저장**하는 Go 애플리케이션입니다.  
Apple의 Bonjour 프로토콜과 호환되며, `grandcat/zeroconf` 라이브러리를 사용해 구현되었습니다.

---

## ✨ 주요 기능

- ✅ UDP 멀티캐스트(mDNS) 기반 데이터 브로드캐스트
- ✅ 데이터 검색 및 수신 후 자동 저장
- ✅ Bonjour/Avahi 호환
- ✅ 인터페이스 지정 없이 로컬 네트워크 전체에 브로드캐스트
- ✅ lightweight 설계, 의존성 최소화

---

## 📦 설치 및 실행 방법

### 1. 설치
```bash
git clone https://github.com/your-org/zeroconf-udp-broadcaster.git
cd zeroconf-udp-broadcaster
go mod tidy
