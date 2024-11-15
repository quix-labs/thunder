const __vite__mapDeps=(i,m=__vite__mapDeps,d=(m.f||(m.f=["./CqY6xSgL.js","./DDeaWyEU.js","./5kaKV4CP.js","./entry.CZH78m3A.css","./JaCKHA-1.js"])))=>i.map(i=>d[i]);
import{aj as z,h as R,m as L,j as K,H as J,av as Q,p as T,o as m,x as C,w as n,b as u,q as s,K as U,c as y,z as b,a as M,t as D,F as N,L as W,d as j,aw as _,aB as X,am as Y,e as Z,f as ee}from"./5kaKV4CP.js";import{z as S,h as se,u as F,a as te,b as oe,c as re,d as ae}from"./JaCKHA-1.js";import{_ as ne}from"./DlAUqK2U.js";const ge=()=>z("sources","/sources"),ie=()=>z("source-drivers","/source-drivers"),ce=Z(()=>ee(()=>import("./CqY6xSgL.js"),__vite__mapDeps([0,1,2,3,4]),import.meta.url).then(d=>d.default||d)),de=["innerHTML"],ue=["textContent"],le={class:"list-disc list-outside ml-4"},me=["textContent"],fe=R({__name:"Source",props:L({mode:{default:"create"},source:{}},{open:{type:Boolean},openModifiers:{}}),emits:L(["created","updated"],["update:open"]),setup(d,{emit:x}){const h=K(d,"open"),w=x,$=J(),{data:f}=ie(),I=Q({driver:null,config:{}}),o=T(()=>d.source||I),A=T(()=>{const e=Object.keys(f.value||{});return S.object({driver:e.length>0?S.enum([e[0],...e.slice(1)],{required_error:"Driver is required"}):S.string({required_error:"Driver is required"})})}),H=async()=>{var e,t,p,l;if(d.mode==="create"){const{status:r,error:i,data:a}=await F("/sources",{method:"post",body:o,watch:!1});if(r.value==="success"){const v=a.value;_().add({title:"Success",color:"success",description:v.message}),w("created")}else r.value==="error"&&_().add({title:"Error",description:((t=(e=i.value)==null?void 0:e.data)==null?void 0:t.error)||((p=i.value)==null?void 0:p.message),color:"error"})}else if(d.mode==="edit"){const{status:r,error:i,data:a}=await F(`/sources/${d.source.id}`,{method:"put",body:o,watch:!1});if(r.value==="success"){const v=a.value;_().add({title:"Success",color:"success",description:v.message}),w("updated")}else r.value==="error"&&_().add({title:"Error",description:(l=i.value)==null?void 0:l.message,color:"error"})}h.value=!1},O=async()=>{var l,r,i;if(await((l=$.value)==null?void 0:l.validate(void 0,{silent:!0}))===!1)return;const{status:e,error:t,data:p}=await F("/source-drivers/test",{method:"post",body:o,watch:!1});if(e.value==="success"){const a=p.value;_().add({title:"Success",color:"success",description:a.message});return}if(e.value==="error"&&[400,422,500].includes(((r=t.value)==null?void 0:r.statusCode)||0)){const a=(i=t.value)==null?void 0:i.data;_().add({title:"Error",description:a==null?void 0:a.error,color:"error"})}};return(e,t)=>{var k,E;const p=te,l=oe,r=re,i=X,a=ce,v=Y,G=ae;return m(),C(G,{title:{create:"New source",edit:`Edit source n°${(k=e.source)==null?void 0:k.id}`,read:`Source n° ${(E=e.source)==null?void 0:E.id}`}[e.mode],description:"Configure data source",open:h.value,"onUpdate:open":t[1]||(t[1]=g=>h.value=g)},{body:n(()=>[u(s(se),{class:"space-y-4",state:s(o),schema:s(A),onSubmit:U(H,["prevent"]),ref_key:"formEl",ref:$,disabled:e.mode==="read",id:"sourceForm"},{default:n(()=>{var g,V,q,B;return[u(l,{name:"driver",label:"Driver",required:""},{default:n(()=>[u(p,{options:Object.entries(s(f)||{}).map(([c,P])=>({value:c,item:P})),modelValue:s(o).driver,"onUpdate:modelValue":t[0]||(t[0]=c=>s(o).driver=c)},{default:n(({item:c})=>[c.config.image?(m(),y("span",{key:0,innerHTML:c.config.image,class:"nested-svg"},null,8,de)):b("",!0),M("span",{textContent:D(c.config.name)},null,8,ue)]),_:1},8,["options","modelValue"])]),_:1}),s(f)&&s(o).driver?(m(),y(N,{key:0},[(((B=(q=(V=(g=s(f))==null?void 0:g[s(o).driver])==null?void 0:V.config)==null?void 0:q.notes)==null?void 0:B.length)||0)>0?(m(),C(i,{key:0,"default-open":""},{content:n(()=>[M("ul",le,[(m(!0),y(N,null,W(s(f)[s(o).driver].config.notes,c=>(m(),y("li",{textContent:D(c)},null,8,me))),256))])]),default:n(()=>[u(r,{label:"Driver Notes",type:"dashed"})]),_:1})):b("",!0),u(i,{"default-open":""},{content:n(()=>[u(a,{state:s(o).config,disabled:e.mode==="read",fields:s(f)[s(o).driver].fields},null,8,["state","disabled","fields"])]),default:n(()=>[u(r,{label:"Driver Configuration"})]),_:1})],64)):b("",!0)]}),_:1},8,["state","schema","disabled"])]),footer:n(()=>[u(v,{type:"button",onClick:U(O,["prevent"]),variant:"soft",color:"secondary"},{default:n(()=>t[2]||(t[2]=[j(" Test configuration ")])),_:1}),e.mode!=="read"?(m(),C(v,{key:0,type:"submit",class:"ml-auto",form:"sourceForm"},{default:n(()=>[j(D({create:"Create",edit:"Save changes"}[e.mode]),1)]),_:1})):b("",!0)]),_:1},8,["title","open"])}}}),ye=ne(fe,[["__scopeId","data-v-83b75755"]]);export{ye as F,ge as u};