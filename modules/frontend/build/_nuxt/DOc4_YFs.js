import{g as D,f as E,r as g,o as N,c as V,b as n,w as i,N as r,V as S,a as m,O as u,d as F,t as M,W as O,Z as y}from"./CD5hUaAV.js";import{_ as R,d as j,c as q,e as A,b as G}from"./CuAb3WZH.js";import{u as L,_ as W}from"./CQ5ZdX0U.js";import"./C_TKe6c3.js";import"./DEi0poqG.js";const Z={class:"flex flex-col flex-1"},H={class:"flex gap-1"},Y=D({__name:"targets",setup(I){const{status:x,data:p,refresh:f}=L(),C=[{key:"id",label:"#",sortable:!0,rowClass:"w-[1px] whitespace-nowrap"},{key:"excerpt",label:"Excerpt",sortable:!0},{key:"driver",label:"Driver",sortable:!0},{key:"actions",sortable:!1,rowClass:"w-[1px] whitespace-nowrap"}],k=E(()=>{var o;return((o=p.value)==null?void 0:o.map(e=>({id:e.id,excerpt:e.excerpt,driver:e.driver})))||[]}),l=g(!1),d=g("create"),c=g(),h=()=>{d.value="create",c.value=void 0,l.value=!0},b=o=>{var e,t;d.value="read",c.value=(t=(e=p.value)==null?void 0:e.filter(s=>s.id===o))==null?void 0:t.at(0),l.value=!0},T=o=>{var e,t;d.value="edit",c.value=(t=(e=p.value)==null?void 0:e.filter(s=>s.id===o))==null?void 0:t.at(0),l.value=!0},U=o=>{var e,t;d.value="create",c.value={...(t=(e=p.value)==null?void 0:e.filter(s=>s.id===o))==null?void 0:t.at(0)},l.value=!0},z=async o=>{var v;const{data:e,error:t,status:s}=await G(`/targets/${o}`,{method:"DELETE"});if(s.value==="error")y().add({color:"red",title:"Unable to delete target",description:(v=t.value)==null?void 0:v.message});else if(s.value==="success"){const _=e.value;y().add({color:"green",title:"Successfully deleted target",description:_.message})}await f()};return(o,e)=>{const t=O,s=j,v=q,_=A,$=R,B=W;return N(),V("section",Z,[n($,{class:"w-full",ui:{divide:"divide-y divide-gray-200 dark:divide-gray-700",body:{padding:""},header:{base:"flex gap-x-2 justify-between items-center"},footer:{base:"text-sm leading-5 text-center"}}},{header:i(()=>[e[4]||(e[4]=m("h1",{class:"font-semibold text-xl text-gray-900 dark:text-white leading-tight"},"Targets",-1)),n(t,{onClick:u(h,["prevent"]),variant:"soft"},{default:i(()=>e[3]||(e[3]=[F("+ Add Target")])),_:1})]),default:i(()=>[n(_,{onSelect:e[1]||(e[1]=a=>b(a.id)),columns:C,rows:r(k),sort:{column:"id",direction:"desc"},loading:r(x)==="pending"||r(x)==="idle"},{"driver-data":i(({row:a})=>[n(s,{size:"xs",label:a.driver,color:"sky",variant:"subtle"},null,8,["label"])]),"actions-data":i(({row:a})=>[m("div",H,[n(v,{items:[[{label:"Replicate",click:()=>U(a.id)}]],onClick:e[0]||(e[0]=u(()=>{},["stop"]))},{default:i(()=>[n(t,{icon:"i-heroicons-ellipsis-horizontal",variant:"link",color:"gray",size:"xl",padded:!1})]),_:2},1032,["items"]),n(t,{icon:"i-heroicons-eye",variant:"link",color:"gray",size:"xl",padded:!1,onClick:u(w=>b(a.id),["stop","prevent"])},null,8,["onClick"]),n(t,{icon:"i-heroicons-pencil-square",variant:"link",color:"gray",size:"xl",padded:!1,onClick:u(w=>T(a.id),["stop","prevent"])},null,8,["onClick"]),n(t,{icon:"i-heroicons-trash",variant:"link",color:"red",size:"xl",padded:!1,onClick:u(w=>z(a.id),["stop","prevent"])},null,8,["onClick"])])]),_:1},8,["rows","loading"])]),footer:i(()=>{var a;return[e[5]||(e[5]=m("span",{class:"font-medium"},"Total: ",-1)),m("span",null,M(((a=r(k))==null?void 0:a.length)||0)+" targets",1)]}),_:1}),n(B,{onUpdated:r(f),onCreated:r(f),mode:r(d),opened:r(l),"onUpdate:opened":e[2]||(e[2]=a=>S(l)?l.value=a:null),target:r(c)},null,8,["onUpdated","onCreated","mode","opened","target"])])}}});export{Y as default};
