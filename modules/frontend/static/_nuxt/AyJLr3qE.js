import{g as A,h as z,i as S,O as N,r as T,G,p as O,k as q,N as M,o as c,x as w,w as l,q as a,aU as X,C as v,z as y,b as m,aV as Y,aW as H,v as h,aX as J,A as E,aY as K,c as F,aZ as W,d as V,t as I,a_ as Z,a$ as Q,am as L,b0 as ee,b1 as ae,b2 as te,D as oe,aj as se,ap as le,av as re,ai as ne,aA as ie,aw as D,as as ce}from"./5kaKV4CP.js";import{h as ue,u as de,b as pe}from"./JaCKHA-1.js";import{_ as me}from"./DDeaWyEU.js";const fe={slots:{overlay:"fixed inset-0 bg-[var(--ui-bg-elevated)]/75",content:"fixed w-full h-dvh bg-[var(--ui-bg)] divide-y divide-[var(--ui-border)] flex flex-col focus:outline-none",header:"px-4 py-5 sm:px-6",body:"flex-1 p-4 sm:p-6",footer:"flex items-center gap-1.5 p-4 sm:px-6",title:"text-[var(--ui-text-highlighted)] font-semibold",description:"mt-1 text-[var(--ui-text-muted)] text-sm",close:"absolute top-4 end-4"},variants:{transition:{true:{overlay:"data-[state=open]:animate-[fade-in_200ms_ease-out] data-[state=closed]:animate-[fade-out_200ms_ease-in]",content:"data-[state=open]:animate-[scale-in_200ms_ease-out] data-[state=closed]:animate-[scale-out_200ms_ease-in]"}},fullscreen:{true:{content:"inset-0"},false:{content:"top-[50%] left-[50%] translate-x-[-50%] translate-y-[-50%] sm:max-w-lg sm:h-auto sm:my-8 sm:rounded-[calc(var(--ui-radius)*2)] sm:shadow-lg sm:ring ring-[var(--ui-border)]"}}}},ve=oe;var P;const ye=A({extend:A(fe),...((P=ve.ui)==null?void 0:P.modal)||{}}),be=z({__name:"Modal",props:{title:{},description:{},content:{},overlay:{type:Boolean,default:!0},transition:{type:Boolean,default:!0},fullscreen:{type:Boolean},portal:{type:Boolean,default:!0},close:{type:Boolean,default:!0},closeIcon:{},preventClose:{type:Boolean},class:{},ui:{},open:{type:Boolean},defaultOpen:{type:Boolean},modal:{type:Boolean,default:!0}},emits:["update:open"],setup(k,{emit:g}){const t=k,$=g,r=S(),R=N(T(t,"open","defaultOpen","modal"),$),x=G(()=>t.content),u=O(()=>t.preventClose?{pointerDownOutside:e=>e.preventDefault(),interactOutside:e=>e.preventDefault(),escapeKeyDown:e=>e.preventDefault()}:{}),B=q(),{t:U}=M(),s=O(()=>ye({transition:t.transition,fullscreen:t.fullscreen}));return(e,n)=>(c(),w(a(te),ee(ae(a(R))),{default:l(({open:f})=>[r.default?(c(),w(a(X),{key:0,"as-child":""},{default:l(()=>[v(e.$slots,"default",{open:f})]),_:2},1024)):y("",!0),m(a(Y),{disabled:!e.portal},{default:l(()=>{var b,_;return[e.overlay?(c(),w(a(H),{key:0,class:h(s.value.overlay({class:(b=t.ui)==null?void 0:b.overlay}))},null,8,["class"])):y("",!0),m(a(J),E({class:s.value.content({class:[t.class,(_=t.ui)==null?void 0:_.content]})},x.value,K(u.value)),{default:l(()=>[v(e.$slots,"content",{},()=>{var i,d,o;return[r.header||e.title||r.title||e.description||r.description||e.close||r.close?(c(),F("div",{key:0,class:h(s.value.header({class:(i=t.ui)==null?void 0:i.header}))},[v(e.$slots,"header",{},()=>{var p,C;return[e.title||r.title?(c(),w(a(W),{key:0,class:h(s.value.title({class:(p=t.ui)==null?void 0:p.title}))},{default:l(()=>[v(e.$slots,"title",{},()=>[V(I(e.title),1)])]),_:3},8,["class"])):y("",!0),e.description||r.description?(c(),w(a(Z),{key:1,class:h(s.value.description({class:(C=t.ui)==null?void 0:C.description}))},{default:l(()=>[v(e.$slots,"description",{},()=>[V(I(e.description),1)])]),_:3},8,["class"])):y("",!0),m(a(Q),{"as-child":""},{default:l(()=>[v(e.$slots,"close",{ui:s.value},()=>{var j;return[e.close?(c(),w(L,E({key:0,icon:e.closeIcon||a(B).ui.icons.close,size:"md",color:"neutral",variant:"ghost","aria-label":a(U)("modal.close")},typeof e.close=="object"?e.close:void 0,{class:s.value.close({class:(j=t.ui)==null?void 0:j.close})}),null,16,["icon","aria-label","class"])):y("",!0)]})]),_:3})]})],2)):y("",!0),r.body?(c(),F("div",{key:1,class:h(s.value.body({class:(d=t.ui)==null?void 0:d.body}))},[v(e.$slots,"body")],2)):y("",!0),r.footer?(c(),F("div",{key:2,class:h(s.value.footer({class:(o=t.ui)==null?void 0:o.footer}))},[v(e.$slots,"footer")],2)):y("",!0)]})]),_:3},16,["class"])]}),_:3},8,["disabled"])]),_:3},16))}}),_e=k=>se("exporters","/exporters",void 0,k),ke=z({__name:"DownloadProcessorForm",props:{processorId:{}},async setup(k){var s;let g,t;const{data:$,status:r,error:R}=([g,t]=le(()=>_e({lazy:!1})),g=await g,t(),g),x=k,u=re({exporter:r.value==="success"&&$.value?(s=Object.keys($.value))==null?void 0:s.at(0):null,filename:null}),B=O(()=>{var e;return`processor-${x.processorId}.${(e=u.exporter)==null?void 0:e.replace("thunder.","")}`}),U=async()=>{var _;const e={...ne(u)};e.filename||(e.filename=B.value);const{data:n,error:f,status:b}=await de(`processors/${x.processorId}/download`,{params:e,responseType:"blob"});if((_=ie())==null||_.close(),b.value==="success"&&n.value){const i=()=>{const d=window.URL.createObjectURL(n.value),o=document.createElement("a");o.href=d,o.download=e.filename,o.click(),setTimeout(()=>{window.URL.revokeObjectURL(d)},0)};i(),D().add({color:"success",title:"File Ready for Download",close:!1,actions:[{icon:"i-heroicons-arrow-down-tray",label:"Download file",color:"primary",variant:"outline",onClick:i}]})}else if(f.value)if(f.value.data instanceof Blob){const i=new FileReader;i.onload=d=>{var o;try{const p=JSON.parse(((o=d.target)==null?void 0:o.result)||"null");D().add({title:"Failed to generate file",description:p.error||"An unknown error occurred.",color:"error"})}catch{D().add({title:"Failed to generate file",description:"Error parsing the error response.",color:"error"})}},i.readAsText(f.value.data)}else D().add({title:"Failed to generate file",description:f.value.message,color:"error"})};return(e,n)=>{const f=ce,b=pe,_=me,i=L,d=be;return c(),w(d,{title:`Download data for processor n°${e.processorId}`},{body:l(()=>[m(a(ue),{state:a(u),onSubmit:U,class:"space-y-4"},{default:l(()=>[m(b,{label:"Exporter",name:"exporter",required:""},{default:l(()=>{var o;return[m(f,{modelValue:a(u).exporter,"onUpdate:modelValue":n[0]||(n[0]=p=>a(u).exporter=p),items:(o=Object.entries(a($)||{}))==null?void 0:o.map(([p,C])=>({label:C,value:p,description:"test"})),required:""},null,8,["modelValue","items"])]}),_:1}),m(b,{label:"Filename",name:"filename"},{default:l(()=>[m(_,{placeholder:a(B),modelValue:a(u).filename,"onUpdate:modelValue":n[1]||(n[1]=o=>a(u).filename=o)},null,8,["placeholder","modelValue"])]),_:1}),m(i,{type:"submit",class:"p-3 cursor-pointer",block:"","loading-auto":""},{default:l(()=>n[2]||(n[2]=[V(" Generate file ")])),_:1})]),_:1},8,["state"])]),_:1},8,["title"])}}});export{ke as default};