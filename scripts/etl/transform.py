import csv
import json
import os
from datetime import datetime

def parse_amount(amount_str):
    try:
        val = float(amount_str)
        return val if val > 0 else None
    except (ValueError, TypeError):
        return None

def standardize_date(date_str):
    if not date_str:
        return datetime.utcnow().isoformat() + "Z"
    
    # Intento de parseo de múltiples formatos inconsistentes
    formats = ["%Y-%m-%dT%H:%M:%SZ", "%Y-%m-%dT%H:%M:%S", "%d/%m/%Y"]
    for fmt in formats:
        try:
            dt = datetime.strptime(date_str.strip(), fmt)
            return dt.isoformat() + "Z"
        except ValueError:
            continue
    return datetime.utcnow().isoformat() + "Z"

def process_etl(input_path, output_path):
    processed_data = []

    with open(input_path, mode='r', encoding='utf-8') as infile:
        reader = csv.DictReader(infile)
        
        for row in reader:
            # 1. Limpieza de valores nulos y estandarización básica
            source = row.get('source_account', '').strip() or 'UNKNOWN'
            target = row.get('target_account', '').strip() or 'UNKNOWN'
            currency = row.get('currency', '').strip().upper()
            
            # 2. Validación de monto
            amount = parse_amount(row.get('amount'))
            if amount is None:
                continue # Descartar transacciones sin monto válido

            # 3. Estandarización de fechas
            timestamp = standardize_date(row.get('timestamp'))

            # 4. Estructuración del payload optimizado para IA
            ai_payload = {
                "features": {
                    "transaction_id": row.get('transaction_id', '').strip(),
                    "source_account": source,
                    "target_account": target,
                    "amount_normalized": amount,
                    "currency": currency,
                    "hour_of_day": datetime.fromisoformat(timestamp.replace('Z', '')).hour
                },
                "metadata": {
                    "original_status": row.get('status', '').strip(),
                    "timestamp": timestamp
                }
            }
            
            processed_data.append(ai_payload)

    # Exportar a JSON estructurado
    with open(output_path, mode='w', encoding='utf-8') as outfile:
        json.dump(processed_data, outfile, indent=2)
    
    print(f"ETL completado. {len(processed_data)} registros procesados y estructurados para IA en {output_path}")

if __name__ == "__main__":
    base_dir = os.path.dirname(os.path.abspath(__file__))
    input_csv = os.path.join(base_dir, 'data', 'raw_transactions.csv')
    output_json = os.path.join(base_dir, 'data', 'ai_ready_transactions.json')
    
    process_etl(input_csv, output_json)
