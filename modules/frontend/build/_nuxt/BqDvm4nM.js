import{_ as M,p as Y,q as k,g as P,x as D,y as V,h as ce,f as u,r as F,L as Q,a1 as G,o as s,c as o,a as O,I as d,a9 as v,ak as B,K as z,F as E,d as N,t as A,J as C,aN as fe,v as T,ab as X,B as x,C as S,H,aO as ge,ah as me,z as pe,ac as ye,k as _,b as K,aP as ve,a7 as be,O as he,aQ as $e}from"./COGry_M7.js";const ze={wrapper:"",inner:"",label:{wrapper:"flex content-center items-center justify-between",base:"block font-medium text-gray-700 dark:text-gray-200",required:"after:content-['*'] after:ms-0.5 after:text-red-500 dark:after:text-red-400"},size:{"2xs":"text-xs",xs:"text-xs",sm:"text-sm",md:"text-sm",lg:"text-sm",xl:"text-base"},container:"mt-1 relative",description:"text-gray-500 dark:text-gray-400",hint:"text-gray-500 dark:text-gray-400",help:"mt-2 text-gray-500 dark:text-gray-400",error:"mt-2 text-red-500 dark:text-red-400",default:{size:"sm"}},ke={wrapper:{base:"flex items-center align-center text-center",horizontal:"w-full flex-row",vertical:"flex-col"},container:{base:"font-medium text-gray-700 dark:text-gray-200 flex",horizontal:"mx-3 whitespace-nowrap",vertical:"my-2"},border:{base:"flex border-gray-200 dark:border-gray-800",horizontal:"w-full",vertical:"h-full",size:{horizontal:{"2xs":"border-t",xs:"border-t-[2px]",sm:"border-t-[3px]",md:"border-t-[4px]",lg:"border-t-[5px]",xl:"border-t-[6px]"},vertical:{"2xs":"border-s",xs:"border-s-[2px]",sm:"border-s-[3px]",md:"border-s-[4px]",lg:"border-s-[5px]",xl:"border-s-[6px]"}},type:{solid:"border-solid",dotted:"border-dotted",dashed:"border-dashed"}},icon:{base:"flex-shrink-0 w-5 h-5"},avatar:{base:"flex-shrink-0",size:"2xs"},label:"text-sm",default:{size:"2xs"}},W=Y(k.ui.strategy,k.ui.formGroup,ze),we=P({inheritAttrs:!1,props:{name:{type:String,default:null},size:{type:String,default:null,validator(e){return Object.keys(W.size).includes(e)}},label:{type:String,default:null},description:{type:String,default:null},required:{type:Boolean,default:!1},help:{type:String,default:null},error:{type:[String,Boolean],default:null},hint:{type:String,default:null},class:{type:[String,Object,Array],default:()=>""},ui:{type:Object,default:()=>({})},eagerValidation:{type:Boolean,default:!1}},setup(e){const{ui:t,attrs:l}=D("formGroup",V(e,"ui"),W,V(e,"class")),a=ce("form-errors",null),c=u(()=>{var m,w;return e.error&&typeof e.error=="string"||typeof e.error=="boolean"?e.error:(w=(m=a==null?void 0:a.value)==null?void 0:m.find(j=>j.path===e.name))==null?void 0:w.message}),r=u(()=>t.value.size[e.size??W.default.size]),g=F(Q());return G("form-group",{error:c,inputId:g,name:u(()=>e.name),size:u(()=>e.size),eagerValidation:u(()=>e.eagerValidation)}),{ui:t,attrs:l,inputId:g,size:r,error:c}}}),Ie=["for"];function Se(e,t,l,a,c,r){return s(),o("div",z({class:e.ui.wrapper},e.attrs),[O("div",{class:d(e.ui.inner)},[e.label||e.$slots.label?(s(),o("div",{key:0,class:d([e.ui.label.wrapper,e.size])},[O("label",{for:e.inputId,class:d([e.ui.label.base,e.required?e.ui.label.required:""])},[e.$slots.label?v(e.$slots,"label",B(z({key:0},{error:e.error,label:e.label,name:e.name,hint:e.hint,description:e.description,help:e.help}))):(s(),o(E,{key:1},[N(A(e.label),1)],64))],10,Ie),e.hint||e.$slots.hint?(s(),o("span",{key:0,class:d([e.ui.hint])},[e.$slots.hint?v(e.$slots,"hint",B(z({key:0},{error:e.error,label:e.label,name:e.name,hint:e.hint,description:e.description,help:e.help}))):(s(),o(E,{key:1},[N(A(e.hint),1)],64))],2)):C("",!0)],2)):C("",!0),e.description||e.$slots.description?(s(),o("p",{key:1,class:d([e.ui.description,e.size])},[e.$slots.description?v(e.$slots,"description",B(z({key:0},{error:e.error,label:e.label,name:e.name,hint:e.hint,description:e.description,help:e.help}))):(s(),o(E,{key:1},[N(A(e.description),1)],64))],2)):C("",!0)],2),O("div",{class:d([e.label?e.ui.container:""])},[v(e.$slots,"default",B(fe({error:e.error}))),typeof e.error=="string"&&e.error?(s(),o("p",{key:0,class:d([e.ui.error,e.size])},[e.$slots.error?v(e.$slots,"error",B(z({key:0},{error:e.error,label:e.label,name:e.name,hint:e.hint,description:e.description,help:e.help}))):(s(),o(E,{key:1},[N(A(e.error),1)],64))],2)):e.help||e.$slots.help?(s(),o("p",{key:1,class:d([e.ui.help,e.size])},[e.$slots.help?v(e.$slots,"help",B(z({key:0},{error:e.error,label:e.label,name:e.name,hint:e.hint,description:e.description,help:e.help}))):(s(),o(E,{key:1},[N(A(e.help),1)],64))],2)):C("",!0)],2)],16)}const De=M(we,[["render",Se]]),U=Y(k.ui.strategy,k.ui.divider,ke),Ce=P({components:{UIcon:T,UAvatar:X},inheritAttrs:!1,props:{label:{type:String,default:null},icon:{type:String,default:null},avatar:{type:Object,default:null},size:{type:String,default:()=>U.default.size,validator(e){return Object.keys(U.border.size.horizontal).includes(e)||Object.keys(U.border.size.vertical).includes(e)}},orientation:{type:String,default:"horizontal",validator:e=>["horizontal","vertical"].includes(e)},type:{type:String,default:"solid",validator:e=>["solid","dotted","dashed"].includes(e)},class:{type:[String,Object,Array],default:()=>""},ui:{type:Object,default:()=>({})}},setup(e){const{ui:t,attrs:l}=D("divider",V(e,"ui"),U),a=u(()=>x(S(t.value.wrapper.base,t.value.wrapper[e.orientation]),e.class)),c=u(()=>S(t.value.container.base,t.value.container[e.orientation])),r=u(()=>S(t.value.border.base,t.value.border[e.orientation],t.value.border.size[e.orientation][e.size],t.value.border.type[e.type]));return{ui:t,attrs:l,wrapperClass:a,containerClass:c,borderClass:r}}});function Oe(e,t,l,a,c,r){const g=T,m=X;return s(),o("div",z({class:e.wrapperClass},e.attrs),[O("div",{class:d(e.borderClass)},null,2),e.label||e.icon||e.avatar||e.$slots.default?(s(),o(E,{key:0},[O("div",{class:d(e.containerClass)},[v(e.$slots,"default",{},()=>[e.label?(s(),o("span",{key:0,class:d(e.ui.label)},A(e.label),3)):e.icon?(s(),H(g,{key:1,name:e.icon,class:d(e.ui.icon.base)},null,8,["name","class"])):e.avatar?(s(),H(m,z({key:2},{size:e.ui.avatar.size,...e.avatar},{class:e.ui.avatar.base}),null,16,["class"])):C("",!0)])],2),O("div",{class:d(e.borderClass)},null,2)],64)):C("",!0)],16)}const Le=M(Ce,[["render",Oe]]),$=Y(k.ui.strategy,k.ui.input,ge),je=P({components:{UIcon:T},inheritAttrs:!1,props:{modelValue:{type:[String,Number],default:""},type:{type:String,default:"text"},id:{type:String,default:null},name:{type:String,default:null},placeholder:{type:String,default:null},required:{type:Boolean,default:!1},disabled:{type:Boolean,default:!1},autofocus:{type:Boolean,default:!1},autofocusDelay:{type:Number,default:100},icon:{type:String,default:null},loadingIcon:{type:String,default:()=>$.default.loadingIcon},leadingIcon:{type:String,default:null},trailingIcon:{type:String,default:null},trailing:{type:Boolean,default:!1},leading:{type:Boolean,default:!1},loading:{type:Boolean,default:!1},padded:{type:Boolean,default:!0},size:{type:String,default:null,validator(e){return Object.keys($.size).includes(e)}},color:{type:String,default:()=>$.default.color,validator(e){return[...k.ui.colors,...Object.keys($.color)].includes(e)}},variant:{type:String,default:()=>$.default.variant,validator(e){return[...Object.keys($.variant),...Object.values($.color).flatMap(t=>Object.keys(t))].includes(e)}},inputClass:{type:String,default:null},class:{type:[String,Object,Array],default:()=>""},ui:{type:Object,default:()=>({})},modelModifiers:{type:Object,default:()=>({})}},emits:["update:modelValue","blur","change"],setup(e,{emit:t,slots:l}){const{ui:a,attrs:c}=D("input",V(e,"ui"),$,V(e,"class")),{size:r,rounded:g}=me({ui:a,props:e}),{emitFormBlur:m,emitFormInput:w,size:j,color:n,inputId:f,name:y}=pe(e,$),p=u(()=>r.value??j.value),b=F(ye({},e.modelModifiers,{trim:!1,lazy:!1,number:!1})),h=F(null),ee=()=>{var i;e.autofocus&&((i=h.value)==null||i.focus())},L=i=>{b.value.trim&&(i=i.trim()),(b.value.number||e.type==="number")&&(i=ve(i)),t("update:modelValue",i),w()},ae=i=>{b.value.lazy||L(i.target.value)},te=i=>{if(e.type==="file"){const I=i.target.files;t("change",I)}else{const I=i.target.value;t("change",I),b.value.lazy&&L(I),b.value.trim&&(i.target.value=I.trim())}},ne=i=>{m(),t("blur",i)};_(()=>{setTimeout(()=>{ee()},e.autofocusDelay)});const re=u(()=>{var I,Z;const i=((Z=(I=a.value.color)==null?void 0:I[n.value])==null?void 0:Z[e.variant])||a.value.variant[e.variant];return x(S(a.value.base,a.value.form,g.value,a.value.placeholder,e.type==="file"&&a.value.file.base,a.value.size[p.value],e.padded?a.value.padding[p.value]:"p-0",i==null?void 0:i.replaceAll("{color}",n.value),(q.value||l.leading)&&a.value.leading.padding[p.value],(R.value||l.trailing)&&a.value.trailing.padding[p.value]),e.inputClass)}),q=u(()=>e.icon&&e.leading||e.icon&&!e.trailing||e.loading&&!e.trailing||e.leadingIcon),R=u(()=>e.icon&&e.trailing||e.loading&&e.trailing||e.trailingIcon),le=u(()=>e.loading?e.loadingIcon:e.leadingIcon||e.icon),ie=u(()=>e.loading&&!q.value?e.loadingIcon:e.trailingIcon||e.icon),se=u(()=>S(a.value.icon.leading.wrapper,a.value.icon.leading.pointer,a.value.icon.leading.padding[p.value])),oe=u(()=>S(a.value.icon.base,n.value&&k.ui.colors.includes(n.value)&&a.value.icon.color.replaceAll("{color}",n.value),a.value.icon.size[p.value],e.loading&&a.value.icon.loading)),ue=u(()=>S(a.value.icon.trailing.wrapper,a.value.icon.trailing.pointer,a.value.icon.trailing.padding[p.value])),de=u(()=>S(a.value.icon.base,n.value&&k.ui.colors.includes(n.value)&&a.value.icon.color.replaceAll("{color}",n.value),a.value.icon.size[p.value],e.loading&&!q.value&&a.value.icon.loading));return{ui:a,attrs:c,name:y,inputId:f,input:h,isLeading:q,isTrailing:R,inputClass:re,leadingIconName:le,leadingIconClass:oe,leadingWrapperIconClass:se,trailingIconName:ie,trailingIconClass:de,trailingWrapperIconClass:ue,onInput:ae,onChange:te,onBlur:ne}}}),Be=["id","name","value","type","required","placeholder","disabled"];function Ee(e,t,l,a,c,r){const g=T;return s(),o("div",{class:d(e.type==="hidden"?"hidden":e.ui.wrapper)},[O("input",z({id:e.inputId,ref:"input",name:e.name,value:e.modelValue,type:e.type,required:e.required,placeholder:e.placeholder,disabled:e.disabled,class:e.inputClass},e.attrs,{onInput:t[0]||(t[0]=(...m)=>e.onInput&&e.onInput(...m)),onBlur:t[1]||(t[1]=(...m)=>e.onBlur&&e.onBlur(...m)),onChange:t[2]||(t[2]=(...m)=>e.onChange&&e.onChange(...m))}),null,16,Be),v(e.$slots,"default"),e.isLeading&&e.leadingIconName||e.$slots.leading?(s(),o("span",{key:0,class:d(e.leadingWrapperIconClass)},[v(e.$slots,"leading",{disabled:e.disabled,loading:e.loading},()=>[K(g,{name:e.leadingIconName,class:d(e.leadingIconClass)},null,8,["name","class"])])],2)):C("",!0),e.isTrailing&&e.trailingIconName||e.$slots.trailing?(s(),o("span",{key:1,class:d(e.trailingWrapperIconClass)},[v(e.$slots,"trailing",{disabled:e.disabled,loading:e.loading},()=>[K(g,{name:e.trailingIconName,class:d(e.trailingIconClass)},null,8,["name","class"])])],2)):C("",!0)],2)}const Re=M(je,[["render",Ee]]);class J extends Error{constructor(t){super(t),this.message=t,Object.setPrototypeOf(this,J.prototype)}}const Ae=P({props:{schema:{type:[Object,Function],default:void 0},state:{type:Object,required:!0},validate:{type:Function,default:()=>[]},validateOn:{type:Array,default:()=>["blur","input","change","submit"]}},emits:["submit","error"],setup(e,{expose:t,emit:l}){const a=Q(),c=$e(`form-${a}`);_(()=>{c.on(async n=>{var f;n.type!=="submit"&&((f=e.validateOn)!=null&&f.includes(n.type))&&await w(n.path,{silent:!0})})}),be(()=>{c.reset()});const r=F([]);G("form-errors",r),G("form-events",c);const g=F({});G("form-inputs",g);async function m(){let n=await e.validate(e.state);if(e.schema)if(qe(e.schema))n=n.concat(await Ue(e.state,e.schema));else if(Ne(e.schema))n=n.concat(await Fe(e.state,e.schema));else if(Ge(e.schema))n=n.concat(await Me(e.state,e.schema));else if(Pe(e.schema))n=n.concat(await Te(e.state,e.schema));else throw new Error("Form validation failed: Unsupported form schema");return n}async function w(n,f={silent:!1}){let y=n;if(n&&!Array.isArray(n)&&(y=[n]),y){const p=r.value.filter(h=>!y.includes(h.path)),b=(await m()).filter(h=>y.includes(h.path));r.value=p.concat(b)}else r.value=await m();if(r.value.length>0){if(f.silent)return!1;throw new J(`Form validation failed: ${JSON.stringify(r.value,null,2)}`)}return e.state}async function j(n){var y;const f=n;try{(y=e.validateOn)!=null&&y.includes("submit")&&await w(),f.data=e.state,l("submit",f)}catch(p){if(!(p instanceof J))throw p;const b={...f,errors:r.value.map(h=>({...h,id:g.value[h.path]}))};l("error",b)}}return t({validate:w,errors:r,setErrors(n,f){f?r.value=r.value.filter(y=>y.path!==f).concat(n):r.value=n},async submit(){await j(new Event("submit"))},getErrors(n){return n?r.value.filter(f=>f.path===n):r.value},clear(n){n?r.value=r.value.filter(f=>f.path!==n):r.value=[]}}),{onSubmit:j}}});function Ne(e){return e.validate&&e.__isYupSchema__}function Ve(e){return e.inner!==void 0}async function Fe(e,t){try{return await t.validate(e,{abortEarly:!1}),[]}catch(l){if(Ve(l))return l.inner.map(a=>({path:a.path??"",message:a.message}));throw l}}function qe(e){return e.parse!==void 0}async function Ue(e,t){const l=await t.safeParseAsync(e);return l.success===!1?l.error.issues.map(a=>({path:a.path.join("."),message:a.message})):[]}function Ge(e){return e.validateAsync!==void 0&&e.id!==void 0}function Je(e){return e.isJoi===!0}async function Me(e,t){try{return await t.validateAsync(e,{abortEarly:!1}),[]}catch(l){if(Je(l))return l.details.map(a=>({path:a.path.join("."),message:a.message}));throw l}}function Pe(e){return"_parse"in e||"_run"in e||typeof e=="function"&&"schema"in e}async function Te(e,t){var a;return((a=(await("_parse"in t?t._parse(e):"_run"in t?t._run({typed:!1,value:e},{}):t(e))).issues)==null?void 0:a.map(c=>{var r;return{path:((r=c.path)==null?void 0:r.map(g=>g.key).join("."))||"",message:c.message}}))||[]}function We(e,t,l,a,c,r){return s(),o("form",{onSubmit:t[0]||(t[0]=he((...g)=>e.onSubmit&&e.onSubmit(...g),["prevent"]))},[v(e.$slots,"default")],32)}const Ze=M(Ae,[["render",We]]);export{De as _,Re as a,Le as b,Ze as c};