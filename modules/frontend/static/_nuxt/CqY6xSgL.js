import{_ as h}from"./DDeaWyEU.js";import{h as _,m as g,bw as x,j as F,p as k,al as U,o as n,x as r,w as p,d,t as f,c as V,F as w,L as B,q as D,as as E}from"./5kaKV4CP.js";import{z as m,b as I,h as M}from"./JaCKHA-1.js";const L=_({name:"DevOnly",inheritAttrs:!1,setup(u,c){return()=>{var t,i;return(i=(t=c.slots).fallback)==null?void 0:i.call(t)}}}),z=_({__name:"DynamicFields",props:g({fields:{default:()=>[]},disabled:{type:Boolean,default:!1}},{state:{required:!0},stateModifiers:{}}),emits:["update:state"],setup(u){const c=x("formEl"),t=F(u,"state"),i=k(()=>{if(u.fields.length===0)return;const l={};for(const e of u.fields){let s={email:m.string({required_error:`${e.label} is required`}).email("Invalid email"),number:m.number({required_error:`${e.label} is required`}),url:m.string({required_error:`${e.label} is required`}).url({message:`${e.label} must be a valid URL`})}[e.type]||m.string({required_error:`${e.label} is required`});e.min&&(s=s.min(parseFloat(e.min),{message:`${e.label} must be at least ${e.min}`})),e.max&&(s=s.max(parseFloat(e.max),{message:`${e.label} must not exceed ${e.max}`})),l[e.name]=e.required?s.min(1,{message:`${e.label} is required`}):m.optional(s).nullable()}return m.object(l)});function b(){t.value||(t.value={}),u.fields.forEach(l=>{var e;l.default&&!((e=t.value)!=null&&e[l.name])&&(t.value[l.name]=l.type!=="number"?l.default:parseFloat(l.default))})}return U(()=>{b()}),(l,e)=>{const s=h,y=E,v=L,q=I,$=M;return n(),r($,{state:t.value,schema:D(i),disabled:l.disabled,ref_key:"formEl",ref:c,class:"space-y-4"},{default:p(()=>[d(f(t.value)+" ",1),(n(!0),V(w,null,B(l.fields,a=>(n(),r(q,{required:a.required,label:a.label,name:`${a.name}`,help:a.help},{default:p(()=>[["number"].includes(a.type)?(n(),r(s,{key:0,class:"w-full",autocomplete:"off",type:"number",min:a.min,max:a.max,value:t.value[a.name],onInput:o=>t.value[a.name]=o.target.value===""?void 0:parseFloat(o.target.value)},null,8,["min","max","value","onInput"])):["url","text","password"].includes(a.type)?(n(),r(s,{key:1,autocomplete:"off",class:"w-full",type:a.type,modelValue:t.value[a.name],"onUpdate:modelValue":o=>t.value[a.name]=o},null,8,["type","modelValue","onUpdate:modelValue"])):a.type==="select"?(n(),r(y,{key:2,autocomplete:"off",class:"w-full",items:a.options,modelValue:t.value[a.name],"onUpdate:modelValue":o=>t.value[a.name]=o},null,8,["items","modelValue","onUpdate:modelValue"])):(n(),r(v,{key:3},{default:p(()=>[d(f(a.type)+" not supported!",1)]),_:2},1024))]),_:2},1032,["required","label","name","help"]))),256))]),_:1},8,["state","schema","disabled"])}}});export{z as default};