from fastapi import FastAPI
from prometheus_api_client import PrometheusConnect
import pandas as pd
import uvicorn

app = FastAPI()

# Cluster içindeki Prometheus servisine bağlan
prom = PrometheusConnect(url="http://prometheus-operated.monitoring.svc.cluster.local:9090", disable_ssl=True)

@app.get("/analiz")
def analyze_risk():
    try:
        # SORGULAR DÜZELTİLDİ: 'sum' eklendi ve 'demo' ismine geçildi
        q_err = 'sum(rate(demo_requests_total{status="500"}[1m]))'
        q_tot = 'sum(rate(demo_requests_total[1m]))'
        
        d_err = prom.custom_query(q_err)
        d_tot = prom.custom_query(q_tot)
        
        val_err = float(d_err[0]['value'][1]) if d_err else 0
        val_tot = float(d_tot[0]['value'][1]) if d_tot else 0.0001
        
        df = pd.DataFrame({'hata': [val_err], 'toplam': [val_tot]})
        oran = (df['hata'][0] / df['toplam'][0]) * 100
        
        # Risk Skoru Hesaplama
        risk = min(100, oran * 1.5) 
        
        karar = "STABIL"
        if risk > 50: karar = "KRITIK - ROLLBACK"
        elif risk > 10: karar = "UYARI"
        
        return {
            "YapayZeka_Analizi": {
                "Hata_Orani": f"%{oran:.2f}",
                "Risk_Skoru": f"{risk:.2f} / 100",
                "Karar": karar
            }
        }
    except Exception as e:
        return {"Hata": str(e)}

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=5000)