import{g as N,f as E,r as f,o as O,c as V,b as o,w as l,N as s,V as S,a as _,O as u,d as F,t as M,Y as R,Z as C}from"./COGry_M7.js";import{a as j,_ as q,b as A}from"./NgrtYLyQ.js";import{_ as G,b as I}from"./ChJoUEp7.js";import{u as L,_ as Y}from"./BgWqHryg.js";import"./BqDvm4nM.js";const Z={class:"flex flex-col flex-1"},H={class:"flex gap-1"},ee=N({__name:"targets",setup(J){const{status:k,data:p,refresh:g}=L(),T=[{key:"id",label:"#",sortable:!0,rowClass:"w-[1px] whitespace-nowrap"},{key:"name",label:"Name",sortable:!0},{key:"driver",label:"Driver",sortable:!0},{key:"actions",sortable:!1,rowClass:"w-[1px] whitespace-nowrap"}],b=E(()=>{var t;return((t=p.value)==null?void 0:t.map((e,n)=>({id:n,name:"TODO name",driver:e.driver})))||[]}),r=f(!1),i=f("create"),d=f(),c=f(),h=()=>{i.value="create",d.value=void 0,c.value=void 0,r.value=!0},w=t=>{var e;i.value="read",d.value=(e=p.value)==null?void 0:e.at(t),c.value=t,r.value=!0},U=t=>{var e;i.value="edit",d.value=(e=p.value)==null?void 0:e.at(t),c.value=t,r.value=!0},z=t=>{var e;i.value="create",d.value={...(e=p.value)==null?void 0:e.at(t)},c.value=void 0,r.value=!0},D=async t=>{var m;const{data:e,error:n,status:v}=await I(`/targets/${t}`,{method:"DELETE"});if(v.value==="error")C().add({color:"red",title:"Unable to delete target",description:(m=n.value)==null?void 0:m.message});else if(v.value==="success"){const x=e.value;C().add({color:"green",title:"Successfully deleted target",description:x.message})}await g()};return(t,e)=>{const n=R,v=j,m=q,x=A,$=G,B=Y;return O(),V("section",Z,[o($,{class:"w-full",ui:{divide:"divide-y divide-gray-200 dark:divide-gray-700",body:{padding:""},header:{base:"flex gap-x-2 justify-between items-center"},footer:{base:"text-sm leading-5 text-center"}}},{header:l(()=>[e[4]||(e[4]=_("h1",{class:"font-semibold text-xl text-gray-900 dark:text-white leading-tight"},"Targets",-1)),o(n,{onClick:u(h,["prevent"]),variant:"soft"},{default:l(()=>e[3]||(e[3]=[F("+ Add Target")])),_:1})]),default:l(()=>[o(x,{onSelect:e[1]||(e[1]=a=>w(a.id)),columns:T,rows:s(b),sort:{column:"id",direction:"desc"},loading:s(k)==="pending"||s(k)==="idle"},{"driver-data":l(({row:a})=>[o(v,{size:"xs",label:a.driver,color:"sky",variant:"subtle"},null,8,["label"])]),"actions-data":l(({row:a})=>[_("div",H,[o(m,{items:[[{label:"Replicate",click:()=>z(a.id)}]],onClick:e[0]||(e[0]=u(()=>{},["stop"]))},{default:l(()=>[o(n,{icon:"i-heroicons-ellipsis-horizontal",variant:"link",color:"gray",size:"xl",padded:!1})]),_:2},1032,["items"]),o(n,{icon:"i-heroicons-eye",variant:"link",color:"gray",size:"xl",padded:!1,onClick:u(y=>w(a.id),["stop","prevent"])},null,8,["onClick"]),o(n,{icon:"i-heroicons-pencil-square",variant:"link",color:"gray",size:"xl",padded:!1,onClick:u(y=>U(a.id),["stop","prevent"])},null,8,["onClick"]),o(n,{icon:"i-heroicons-trash",variant:"link",color:"red",size:"xl",padded:!1,onClick:u(y=>D(a.id),["stop","prevent"])},null,8,["onClick"])])]),_:1},8,["rows","loading"])]),footer:l(()=>{var a;return[e[5]||(e[5]=_("span",{class:"font-medium"},"Total: ",-1)),_("span",null,M(((a=s(b))==null?void 0:a.length)||0)+" targets",1)]}),_:1}),o(B,{onUpdated:s(g),onCreated:s(g),mode:s(i),opened:s(r),"onUpdate:opened":e[2]||(e[2]=a=>S(r)?r.value=a:null),target:s(d),"processor-id":s(c)},null,8,["onUpdated","onCreated","mode","opened","target","processor-id"])])}}});export{ee as default};