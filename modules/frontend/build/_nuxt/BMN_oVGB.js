import{g as E,f as T,r as _,o as N,c as V,b as a,w as l,N as n,V as F,a as f,O as u,d as O,t as j,Y as M,Z as C}from"./BtwgJjUD.js";import{_ as R,d as q,c as A,e as G,b as I}from"./CElnJIvZ.js";import{u as L,_ as Y}from"./Txo80PP-.js";import"./DBOZ3qfk.js";const Z={class:"flex flex-col flex-1"},H={class:"flex gap-1"},X=E({__name:"sources",setup(J){const{status:k,data:p,refresh:x}=L(),h=[{key:"id",label:"#",sortable:!0,rowClass:"w-[1px] whitespace-nowrap"},{key:"excerpt",label:"Excerpt",sortable:!0},{key:"driver",label:"Driver",sortable:!0},{key:"actions",sortable:!1,rowClass:"w-[1px] whitespace-nowrap"}],b=T(()=>{var o;return((o=Object.entries(p.value||{}))==null?void 0:o.map(([e,s])=>({id:parseFloat(e),excerpt:s.excerpt,driver:s.driver})))||[]}),r=_(!1),i=_("create"),d=_(),c=_(),S=()=>{i.value="create",d.value=void 0,c.value=void 0,r.value=!0},w=o=>{var e;i.value="read",d.value=(e=p.value)==null?void 0:e[o],c.value=o,r.value=!0},U=o=>{var e;i.value="edit",d.value=(e=p.value)==null?void 0:e[o],c.value=o,r.value=!0},z=o=>{var e;i.value="create",d.value={...(e=p.value)==null?void 0:e[o]},c.value=void 0,r.value=!0},$=async o=>{var m;const{data:e,error:s,status:v}=await I(`/sources/${o}`,{method:"DELETE"});if(v.value==="error")C().add({color:"red",title:"Unable to delete source",description:(m=s.value)==null?void 0:m.message});else if(v.value==="success"){const g=e.value;C().add({color:"green",title:"Successfully deleted source",description:g.message})}await x()};return(o,e)=>{const s=M,v=q,m=A,g=G,B=R,D=Y;return N(),V("section",Z,[a(B,{class:"w-full",ui:{divide:"divide-y divide-gray-200 dark:divide-gray-700",body:{padding:""},header:{base:"flex gap-x-2 justify-between items-center"},footer:{base:"text-sm leading-5 text-center"}}},{header:l(()=>[e[4]||(e[4]=f("h1",{class:"font-semibold text-xl text-gray-900 dark:text-white leading-tight"},"Sources",-1)),a(s,{onClick:u(S,["prevent"]),variant:"soft"},{default:l(()=>e[3]||(e[3]=[O("+ Add Source")])),_:1})]),default:l(()=>[a(g,{onSelect:e[1]||(e[1]=t=>w(t.id)),columns:h,rows:n(b),sort:{column:"id",direction:"desc"},loading:n(k)==="pending"||n(k)==="idle"},{"driver-data":l(({row:t})=>[a(v,{size:"xs",label:t.driver,color:"sky",variant:"subtle"},null,8,["label"])]),"actions-data":l(({row:t})=>[f("div",H,[a(m,{items:[[{label:"Replicate",click:()=>z(t.id)}]],onClick:e[0]||(e[0]=u(()=>{},["stop"]))},{default:l(()=>[a(s,{icon:"i-heroicons-ellipsis-horizontal",variant:"link",color:"gray",size:"xl",padded:!1})]),_:2},1032,["items"]),a(s,{icon:"i-heroicons-eye",variant:"link",color:"gray",size:"xl",padded:!1,onClick:u(y=>w(t.id),["stop","prevent"])},null,8,["onClick"]),a(s,{icon:"i-heroicons-pencil-square",variant:"link",color:"gray",size:"xl",padded:!1,onClick:u(y=>U(t.id),["stop","prevent"])},null,8,["onClick"]),a(s,{icon:"i-heroicons-trash",variant:"link",color:"red",size:"xl",padded:!1,onClick:u(y=>$(t.id),["stop","prevent"])},null,8,["onClick"])])]),_:1},8,["rows","loading"])]),footer:l(()=>{var t;return[e[5]||(e[5]=f("span",{class:"font-medium"},"Total: ",-1)),f("span",null,j(((t=n(b))==null?void 0:t.length)||0)+" sources",1)]}),_:1}),a(D,{onUpdated:n(x),onCreated:n(x),mode:n(i),opened:n(r),"onUpdate:opened":e[2]||(e[2]=t=>F(r)?r.value=t:null),source:n(d),"source-id":n(c)},null,8,["onUpdated","onCreated","mode","opened","source","source-id"])])}}});export{X as default};
