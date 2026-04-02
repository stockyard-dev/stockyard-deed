package server
import("encoding/json";"log";"net/http";"github.com/stockyard-dev/stockyard-deed/internal/store")
type Server struct{db *store.DB;mux *http.ServeMux}
func New(db *store.DB)*Server{s:=&Server{db:db,mux:http.NewServeMux()}
s.mux.HandleFunc("GET /api/licenses",s.list);s.mux.HandleFunc("POST /api/licenses",s.create);s.mux.HandleFunc("GET /api/licenses/{id}",s.get);s.mux.HandleFunc("DELETE /api/licenses/{id}",s.del)
s.mux.HandleFunc("POST /api/licenses/{id}/revoke",s.revoke);s.mux.HandleFunc("GET /api/validate",s.validate)
s.mux.HandleFunc("GET /api/stats",s.stats);s.mux.HandleFunc("GET /api/health",s.health)
s.mux.HandleFunc("GET /ui",s.dashboard);s.mux.HandleFunc("GET /ui/",s.dashboard);s.mux.HandleFunc("GET /",s.root);return s}
func(s *Server)ServeHTTP(w http.ResponseWriter,r *http.Request){s.mux.ServeHTTP(w,r)}
func wj(w http.ResponseWriter,c int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(c);json.NewEncoder(w).Encode(v)}
func we(w http.ResponseWriter,c int,m string){wj(w,c,map[string]string{"error":m})}
func(s *Server)root(w http.ResponseWriter,r *http.Request){if r.URL.Path!="/"{http.NotFound(w,r);return};http.Redirect(w,r,"/ui",302)}
func(s *Server)list(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"licenses":oe(s.db.List())})}
func(s *Server)create(w http.ResponseWriter,r *http.Request){var l store.License;json.NewDecoder(r.Body).Decode(&l);if l.Product==""{we(w,400,"product required");return};s.db.Create(&l);wj(w,201,s.db.Get(l.ID))}
func(s *Server)get(w http.ResponseWriter,r *http.Request){l:=s.db.Get(r.PathValue("id"));if l==nil{we(w,404,"not found");return};wj(w,200,l)}
func(s *Server)del(w http.ResponseWriter,r *http.Request){s.db.Delete(r.PathValue("id"));wj(w,200,map[string]string{"deleted":"ok"})}
func(s *Server)revoke(w http.ResponseWriter,r *http.Request){s.db.Revoke(r.PathValue("id"));wj(w,200,s.db.Get(r.PathValue("id")))}
func(s *Server)validate(w http.ResponseWriter,r *http.Request){key:=r.URL.Query().Get("key");l:=s.db.Validate(key);if l==nil{wj(w,200,map[string]any{"valid":false});return};wj(w,200,map[string]any{"valid":l.Status=="active","license":l})}
func(s *Server)stats(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"total":s.db.Count(),"active":s.db.ActiveCount()})}
func(s *Server)health(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"status":"ok","service":"deed","licenses":s.db.Count()})}
func oe[T any](s []T)[]T{if s==nil{return[]T{}};return s}
func init(){log.SetFlags(log.LstdFlags|log.Lshortfile)}
