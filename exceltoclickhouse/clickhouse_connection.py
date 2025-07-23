import clickhouse_connect

if __name__ == '__main__':
    client = clickhouse_connect.get_client(
        host='npomobbg93.germanywestcentral.azure.clickhouse.cloud',
        user='default',
        password='1S.6V_z9Lr9Wc',
        secure=True
    )
    print("Result:", client.query("SELECT 1").result_set[0][0])
