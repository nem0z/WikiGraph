from utils import get_mysql_conn, get_redis_conn

conn = None
cursor = None
r = None

try:
    conn = get_mysql_conn()
    cursor = conn.cursor()
    cursor.execute("SELECT id, link FROM articles")
    articles = cursor.fetchall()
    if not articles:
        raise Exception("No article to load from MySQL")
    
    r = get_redis_conn()
    for id, link in articles:
        r.set(id, link) # type: ignore
    
    print("Redis population complete")
finally:
    if cursor is not None:
        cursor.close()
    if conn is not None:
        conn.close()
    if r is not None:
        r.close()