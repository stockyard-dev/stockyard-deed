package store
import ("crypto/rand";"database/sql";"encoding/hex";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type License struct{ID string `json:"id"`;Product string `json:"product"`;Key string `json:"key"`;Holder string `json:"holder,omitempty"`;Email string `json:"email,omitempty"`;Status string `json:"status"`;MaxSeats int `json:"max_seats"`;ExpiresAt string `json:"expires_at,omitempty"`;CreatedAt string `json:"created_at"`}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"deed.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS licenses(id TEXT PRIMARY KEY,product TEXT NOT NULL,key TEXT UNIQUE NOT NULL,holder TEXT DEFAULT '',email TEXT DEFAULT '',status TEXT DEFAULT 'active',max_seats INTEGER DEFAULT 1,expires_at TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func genKey()string{b:=make([]byte,16);rand.Read(b);return hex.EncodeToString(b)}
func(d *DB)Create(l *License)error{l.ID=genID();l.CreatedAt=now();if l.Key==""{l.Key=genKey()};if l.Status==""{l.Status="active"};if l.MaxSeats<=0{l.MaxSeats=1}
_,err:=d.db.Exec(`INSERT INTO licenses VALUES(?,?,?,?,?,?,?,?,?)`,l.ID,l.Product,l.Key,l.Holder,l.Email,l.Status,l.MaxSeats,l.ExpiresAt,l.CreatedAt);return err}
func(d *DB)Get(id string)*License{var l License;if d.db.QueryRow(`SELECT * FROM licenses WHERE id=?`,id).Scan(&l.ID,&l.Product,&l.Key,&l.Holder,&l.Email,&l.Status,&l.MaxSeats,&l.ExpiresAt,&l.CreatedAt)!=nil{return nil};return &l}
func(d *DB)Validate(key string)*License{var l License;if d.db.QueryRow(`SELECT * FROM licenses WHERE key=?`,key).Scan(&l.ID,&l.Product,&l.Key,&l.Holder,&l.Email,&l.Status,&l.MaxSeats,&l.ExpiresAt,&l.CreatedAt)!=nil{return nil};return &l}
func(d *DB)List()[]License{rows,_:=d.db.Query(`SELECT * FROM licenses ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close()
var o []License;for rows.Next(){var l License;rows.Scan(&l.ID,&l.Product,&l.Key,&l.Holder,&l.Email,&l.Status,&l.MaxSeats,&l.ExpiresAt,&l.CreatedAt);o=append(o,l)};return o}
func(d *DB)Revoke(id string)error{_,err:=d.db.Exec(`UPDATE licenses SET status='revoked' WHERE id=?`,id);return err}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM licenses WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM licenses`).Scan(&n);return n}
func(d *DB)ActiveCount()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM licenses WHERE status='active'`).Scan(&n);return n}
