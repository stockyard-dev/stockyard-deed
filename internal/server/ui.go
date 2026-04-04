package server
import "net/http"
func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) { w.Header().Set("Content-Type", "text/html"); w.Write([]byte(dashHTML)) }
const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Deed</title><link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet"><style>:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--blue:#5b8dd9;--mono:'JetBrains Mono',monospace}*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5}.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}.hdr h1{font-size:.9rem;letter-spacing:2px}.hdr h1 span{color:var(--rust)}.main{padding:1.5rem;max-width:960px;margin:0 auto}.stats{display:grid;grid-template-columns:repeat(3,1fr);gap:.5rem;margin-bottom:1rem}.st{background:var(--bg2);border:1px solid var(--bg3);padding:.6rem;text-align:center}.st-v{font-size:1.2rem;font-weight:700}.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.15rem}.toolbar{display:flex;gap:.5rem;margin-bottom:1rem}.search{flex:1;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}.search:focus{outline:none;border-color:var(--leather)}.lic{background:var(--bg2);border:1px solid var(--bg3);padding:.8rem 1rem;margin-bottom:.5rem;transition:border-color .2s}.lic:hover{border-color:var(--leather)}.lic-top{display:flex;justify-content:space-between;align-items:flex-start;gap:.5rem}.lic-product{font-size:.85rem;font-weight:700;color:var(--gold)}.lic-holder{font-size:.7rem;color:var(--cd);margin-top:.1rem}.lic-key{font-size:.55rem;color:var(--cm);margin-top:.3rem;background:var(--bg);padding:.2rem .4rem;border:1px solid var(--bg3);font-family:var(--mono);word-break:break-all}.lic-meta{font-size:.55rem;color:var(--cm);margin-top:.3rem;display:flex;gap:.5rem;flex-wrap:wrap}.lic-actions{display:flex;gap:.3rem;flex-shrink:0}.badge{font-size:.5rem;padding:.12rem .35rem;text-transform:uppercase;letter-spacing:1px;border:1px solid}.badge.active{border-color:var(--green);color:var(--green)}.badge.expired{border-color:var(--red);color:var(--red)}.badge.revoked{border-color:var(--cm);color:var(--cm)}.badge.trial{border-color:var(--gold);color:var(--gold)}.btn{font-size:.6rem;padding:.25rem .5rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:all .2s}.btn:hover{border-color:var(--leather);color:var(--cream)}.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}.btn-sm{font-size:.55rem;padding:.2rem .4rem}.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:460px;max-width:92vw}.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust)}.fr{margin-bottom:.6rem}.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}.fr input,.fr select{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}.fr input:focus,.fr select:focus{outline:none;border-color:var(--leather)}.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.75rem}</style></head><body>
<div class="hdr"><h1><span>&#9670;</span> DEED</h1><button class="btn btn-p" onclick="openForm()">+ Issue License</button></div>
<div class="main"><div class="stats" id="stats"></div><div class="toolbar"><input class="search" id="search" placeholder="Search licenses..." oninput="render()"></div><div id="list"></div></div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()"><div class="modal" id="mdl"></div></div>
<script>
var A='/api',items=[],editId=null;
async function load(){var r=await fetch(A+'/licenses').then(function(r){return r.json()});items=r.licenses||[];renderStats();render();}
function renderStats(){var t=items.length,act=items.filter(function(l){return l.status==='active'}).length,seats=items.reduce(function(s,l){return s+(l.max_seats||0)},0);
document.getElementById('stats').innerHTML='<div class="st"><div class="st-v">'+t+'</div><div class="st-l">Licenses</div></div><div class="st"><div class="st-v" style="color:var(--green)">'+act+'</div><div class="st-l">Active</div></div><div class="st"><div class="st-v">'+seats+'</div><div class="st-l">Total Seats</div></div>';}
function render(){var q=(document.getElementById('search').value||'').toLowerCase();var f=items;
if(q)f=f.filter(function(l){return(l.product||'').toLowerCase().includes(q)||(l.holder||'').toLowerCase().includes(q)||(l.email||'').toLowerCase().includes(q)||(l.key||'').toLowerCase().includes(q)});
if(!f.length){document.getElementById('list').innerHTML='<div class="empty">No licenses issued.</div>';return;}
var h='';f.forEach(function(l){
h+='<div class="lic"><div class="lic-top"><div style="flex:1"><div class="lic-product">'+esc(l.product)+'</div>';
if(l.holder||l.email)h+='<div class="lic-holder">'+esc(l.holder||'')+(l.email?' &lt;'+esc(l.email)+'&gt;':'')+'</div>';
h+='</div><div class="lic-actions"><button class="btn btn-sm" onclick="openEdit(''+l.id+'')">Edit</button><button class="btn btn-sm" onclick="del(''+l.id+'')" style="color:var(--red)">&#10005;</button></div></div>';
if(l.key)h+='<div class="lic-key">'+esc(l.key)+'</div>';
h+='<div class="lic-meta">';
if(l.status)h+='<span class="badge '+l.status+'">'+l.status+'</span>';
if(l.max_seats)h+='<span>'+l.max_seats+' seats</span>';
if(l.expires_at)h+='<span>Expires: '+l.expires_at+'</span>';
h+='<span>'+ft(l.created_at)+'</span></div></div>';});
document.getElementById('list').innerHTML=h;}
async function del(id){if(!confirm('Revoke?'))return;await fetch(A+'/licenses/'+id,{method:'DELETE'});load();}
function formHTML(lic){var i=lic||{product:'',key:'',holder:'',email:'',status:'active',max_seats:1,expires_at:''};var isEdit=!!lic;
var h='<h2>'+(isEdit?'EDIT':'ISSUE')+' LICENSE</h2>';
h+='<div class="row2"><div class="fr"><label>Product *</label><input id="f-prod" value="'+esc(i.product)+'"></div><div class="fr"><label>Key</label><input id="f-key" value="'+esc(i.key)+'" placeholder="auto or manual"></div></div>';
h+='<div class="row2"><div class="fr"><label>Holder</label><input id="f-holder" value="'+esc(i.holder)+'"></div><div class="fr"><label>Email</label><input id="f-email" value="'+esc(i.email)+'"></div></div>';
h+='<div class="row2"><div class="fr"><label>Status</label><select id="f-status">';
['active','trial','expired','revoked'].forEach(function(s){h+='<option value="'+s+'"'+(i.status===s?' selected':'')+'>'+s.charAt(0).toUpperCase()+s.slice(1)+'</option>';});
h+='</select></div><div class="fr"><label>Max Seats</label><input id="f-seats" type="number" value="'+(i.max_seats||1)+'"></div></div>';
h+='<div class="fr"><label>Expires</label><input id="f-exp" type="date" value="'+esc(i.expires_at)+'"></div>';
h+='<div class="acts"><button class="btn" onclick="closeModal()">Cancel</button><button class="btn btn-p" onclick="submit()">'+(isEdit?'Save':'Issue')+'</button></div>';return h;}
function openForm(){editId=null;document.getElementById('mdl').innerHTML=formHTML();document.getElementById('mbg').classList.add('open');}
function openEdit(id){var l=null;for(var j=0;j<items.length;j++){if(items[j].id===id){l=items[j];break;}}if(!l)return;editId=id;document.getElementById('mdl').innerHTML=formHTML(l);document.getElementById('mbg').classList.add('open');}
function closeModal(){document.getElementById('mbg').classList.remove('open');editId=null;}
async function submit(){var prod=document.getElementById('f-prod').value.trim();if(!prod){alert('Product required');return;}
var body={product:prod,key:document.getElementById('f-key').value.trim(),holder:document.getElementById('f-holder').value.trim(),email:document.getElementById('f-email').value.trim(),status:document.getElementById('f-status').value,max_seats:parseInt(document.getElementById('f-seats').value)||1,expires_at:document.getElementById('f-exp').value};
if(editId){await fetch(A+'/licenses/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
else{await fetch(A+'/licenses',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}closeModal();load();}
function ft(t){if(!t)return'';try{return new Date(t).toLocaleDateString('en-US',{month:'short',day:'numeric',year:'numeric'})}catch(e){return t;}}
function esc(s){if(!s)return'';var d=document.createElement('div');d.textContent=s;return d.innerHTML;}
document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal();});load();
</script></body></html>`
