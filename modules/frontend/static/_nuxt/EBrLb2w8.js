var X=Object.defineProperty;var ee=(r,e,t)=>e in r?X(r,e,{enumerable:!0,configurable:!0,writable:!0,value:t}):r[e]=t;var q=(r,e,t)=>ee(r,typeof e!="symbol"?e+"":e,t);import{e as L,B as re,h as A,bo as Q,N as D,al as G,bp as S,bq as te,br as ae,o as g,c as j,a as Y,Q as p,b as se,w as H,y as w,d as P,t as C,j as J,ao as ne,m as N,z as W,af as R,bs as M,O as T,bt as U,as as oe,bu as ie,bv as le,bw as ue,bx as ce,i as de,k as pe,b9 as fe,by as me}from"./CmCqKDSp.js";const ve={slots:{root:"",wrapper:"",labelWrapper:"flex content-center items-center justify-between",label:"block font-medium text-[var(--ui-text)]",container:"mt-1 relative",description:"text-[var(--ui-text-muted)]",error:"mt-2 text-[var(--ui-error)]",hint:"text-[var(--ui-text-muted)]",help:"mt-2 text-[var(--ui-text-muted)]"},variants:{size:{xs:{root:"text-xs"},sm:{root:"text-xs"},md:{root:"text-sm"},lg:{root:"text-sm"},xl:{root:"text-base"}},required:{true:{label:"after:content-['*'] after:ms-0.5 after:text-[var(--ui-error)]"}}},defaultVariants:{size:"md"}},ye=R;var Z;const he=W({extend:W(ve),...((Z=ye.ui)==null?void 0:Z.formField)||{}}),qe=L({__name:"FormField",props:{name:{},errorPattern:{},label:{},description:{},help:{},error:{type:[String,Boolean]},hint:{},size:{},required:{type:Boolean},eagerValidation:{type:Boolean},validateOnInputDelay:{},class:{},ui:{}},setup(r){const e=r,t=re(),a=A(()=>he({size:e.size,required:e.required})),i=Q("form-errors",null),l=A(()=>{var n,f;return e.error||((f=(n=i==null?void 0:i.value)==null?void 0:n.find(u=>u.name===e.name||e.errorPattern&&u.name.match(e.errorPattern)))==null?void 0:f.message)}),c=D(G());return S(te,c),S(ae,A(()=>({error:l.value,name:e.name,size:e.size,eagerValidation:e.eagerValidation,validateOnInputDelay:e.validateOnInputDelay,errorPattern:e.errorPattern}))),(n,f)=>{var u,m,I,E,v,y,k,O,s;return g(),j("div",{class:p(a.value.root({class:[e.class,(u=e.ui)==null?void 0:u.root]}))},[Y("div",{class:p(a.value.wrapper({class:(m=e.ui)==null?void 0:m.wrapper}))},[n.label||t.label?(g(),j("div",{key:0,class:p(a.value.labelWrapper({class:(I=e.ui)==null?void 0:I.labelWrapper}))},[se(J(ne),{for:c.value,class:p(a.value.label({class:(E=e.ui)==null?void 0:E.label}))},{default:H(()=>[w(n.$slots,"label",{label:n.label},()=>[P(C(n.label),1)])]),_:3},8,["for","class"]),n.hint||t.hint?(g(),j("span",{key:0,class:p(a.value.hint({class:(v=e.ui)==null?void 0:v.hint}))},[w(n.$slots,"hint",{hint:n.hint},()=>[P(C(n.hint),1)])],2)):N("",!0)],2)):N("",!0),n.description||t.description?(g(),j("p",{key:1,class:p(a.value.description({class:(y=e.ui)==null?void 0:y.description}))},[w(n.$slots,"description",{description:n.description},()=>[P(C(n.description),1)])],2)):N("",!0)],2),Y("div",{class:p([(n.label||!!t.label||n.description||!!t.description)&&a.value.container({class:(k=e.ui)==null?void 0:k.container})])},[w(n.$slots,"default",{error:l.value}),typeof l.value=="string"&&l.value||t.error?(g(),j("p",{key:0,class:p(a.value.error({class:(O=e.ui)==null?void 0:O.error}))},[w(n.$slots,"error",{error:l.value},()=>[P(C(l.value),1)])],2)):n.help||t.help?(g(),j("p",{key:1,class:p(a.value.help({class:(s=e.ui)==null?void 0:s.help}))},[w(n.$slots,"help",{help:n.help},()=>[P(C(n.help),1)])],2)):N("",!0)],2)],2)}}});function be(r){return r.validate&&r.__isYupSchema__}function ge(r){return r.inner!==void 0}function we(r){return"schema"in r&&typeof r.coercer=="function"&&typeof r.validator=="function"&&typeof r.refiner=="function"}function Se(r){return r.parse!==void 0}function Ie(r){return r.validateAsync!==void 0&&r.id!==void 0}function je(r){return r.isJoi===!0}function Ee(r){return"_run"in r||typeof r=="function"&&"schema"in r}function ke(r){return"~standard"in r}async function Oe(r,e){var a;const t=await e["~standard"].validate({value:r});return t.issues?{errors:((a=t.issues)==null?void 0:a.map(i=>{var l;return{name:((l=i.path)==null?void 0:l.map(c=>typeof c=="object"?c.key:c).join("."))||"",message:i.message}}))||[],result:null}:{errors:null,result:t.value}}async function Be(r,e){try{return{errors:null,result:e.validateSync(r,{abortEarly:!1})}}catch(t){if(ge(t))return{errors:t.inner.map(i=>({name:i.path??"",message:i.message})),result:null};throw t}}async function $e(r,e){const[t,a]=e.validate(r);return t?{errors:t.failures().map(l=>({message:l.message,name:l.path.join(".")})),result:null}:{errors:null,result:a}}async function xe(r,e){const t=await e.safeParseAsync(r);return t.success===!1?{errors:t.error.issues.map(i=>({name:i.path.join("."),message:i.message})),result:null}:{result:t.data,errors:null}}async function Ve(r,e){try{return{errors:null,result:await e.validateAsync(r,{abortEarly:!1})}}catch(t){if(je(t))return{errors:t.details.map(i=>({name:i.path.join("."),message:i.message})),result:null};throw t}}async function ze(r,e){const t=await("_run"in e?e._run({typed:!1,value:r},{}):e(r));return!t.issues||t.issues.length===0?{errors:null,result:"output"in t?t.output:"value"in t?t.value:null}:{errors:t.issues.map(i=>{var l;return{name:((l=i.path)==null?void 0:l.map(c=>c.key).join("."))||"",message:i.message}}),result:null}}function Fe(r,e){if(Se(e))return xe(r,e);if(Ie(e))return Ve(r,e);if(Ee(e))return ze(r,e);if(be(e))return Be(r,e);if(we(e))return $e(r,e);if(ke(e))return Oe(r,e);throw new Error("Form validation failed: Unsupported form schema")}class K extends Error{constructor(t,a,i){super("Form validation exception");q(this,"formId");q(this,"errors");q(this,"childrens");this.formId=t,this.errors=a,this.childrens=i,Object.setPrototypeOf(this,K.prototype)}}const Pe={base:""},Ce=R;var _;const De=W({extend:W(Pe),...((_=Ce.ui)==null?void 0:_.form)||{}}),Ne=L({__name:"Form",props:{id:{},schema:{},state:{},validate:{},validateOn:{default(){return["input","blur","change"]}},disabled:{type:Boolean},validateOnInputDelay:{default:300},class:{},onSubmit:{}},emits:["submit","error"],setup(r,{expose:e,emit:t}){const a=r,i=t,l=a.id??G(),c=me(`form-${l}`),n=Q(M,void 0);S(M,c);const f=D(new Map);T(async()=>{c.on(async s=>{var o;s.type==="attach"?f.value.set(s.formId,{validate:s.validate}):s.type==="detach"?f.value.delete(s.formId):(o=a.validateOn)!=null&&o.includes(s.type)&&await v({name:s.name,silent:!0,nested:!1})})}),U(()=>{c.reset()}),T(async()=>{n&&(await oe(),n.emit({type:"attach",validate:v,formId:l}))}),U(()=>{n&&n.emit({type:"detach",formId:l})});const u=D([]);S("form-errors",u);const m=D({});S(ie,m);function I(s){return s.map(o=>{var d;return{...o,id:(d=m.value[o.name])==null?void 0:d.id}})}async function E(){let s=a.validate?await a.validate(a.state)??[]:[];if(a.schema){const{errors:o,result:d}=await Fe(a.state,a.schema);o?s=s.concat(o):Object.assign(a.state,d)}return I(s)}async function v(s={silent:!1,nested:!0}){const o=s.name&&!Array.isArray(s.name)?[s.name]:s.name,d=!o&&s.nested?Array.from(f.value.values()).map(({validate:b})=>b().then(()=>{}).catch(B=>{if(!(B instanceof K))throw B;return B})):[];if(o){const b=u.value.filter($=>!o.some(x=>{var z,F;const V=(F=(z=m.value)==null?void 0:z[x])==null?void 0:F.pattern;return x===$.name||V&&$.name.match(V)})),B=(await E()).filter($=>o.some(x=>{var z,F;const V=(F=(z=m.value)==null?void 0:z[x])==null?void 0:F.pattern;return x===$.name||V&&$.name.match(V)}));u.value=b.concat(B)}else u.value=await E();const h=(await Promise.all(d)).filter(b=>b);if(u.value.length+h.length>0){if(s.silent)return!1;throw new K(l,u.value,h)}return a.state}const y=D(!1);S(ue,le(y));async function k(s){var d;y.value=!0;const o=s;try{await v({nested:!0}),o.data=a.state,await((d=a.onSubmit)==null?void 0:d.call(a,o))}catch(h){if(!(h instanceof K))throw h;const b={...o,errors:h.errors,childrens:h.childrens};i("error",b)}y.value=!1}const O=A(()=>a.disabled||y.value);return S(ce,A(()=>({disabled:O.value,validateOnInputDelay:a.validateOnInputDelay}))),e({validate:v,errors:u,setErrors(s,o){o?u.value=u.value.filter(d=>d.name!==o).concat(I(s)):u.value=I(s)},async submit(){await k(new Event("submit"))},getErrors(s){return s?u.value.filter(o=>o.name===s):u.value},clear(s){s?u.value=u.value.filter(o=>o.name!==s):u.value=[]},disabled:O}),(s,o)=>(g(),de(fe(J(n)?"div":"form"),{id:J(l),class:p(J(De)({class:a.class})),onSubmit:pe(k,["prevent"])},{default:H(()=>[w(s.$slots,"default")]),_:3},40,["id","class"]))}});export{Ne as _,qe as a};