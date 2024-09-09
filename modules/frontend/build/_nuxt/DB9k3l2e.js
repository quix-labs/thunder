import{ai as he,s as fe,d as D,b as T,g as le,am as me,r as ye,M as ce,v as ge,ak as ve,aM as be,o as F,c as M,w as A,a as x,h as v,j as X,a5 as we,t as ee,m as ke,l as Ce,$ as Ne}from"./CcBA2A6A.js";import{b as je,h as Se}from"./n-PNidO-.js";import{_ as _e,u as Oe}from"./B80jb0T7.js";function xe(e,t){let n,r=null;return async function(...i){n?r=i:(n=!0,await e(...i),setTimeout(async()=>{n=!1,r&&(await e(...r),r=null)},t))}}function Ve(e,t="data.json"){const n=JSON.stringify(e,null,2),r=new Blob([n],{type:"application/json"}),i=new File([r],t,{type:"application/json"}),p=URL.createObjectURL(i),f=document.createElement("a");f.href=p,f.download=t,document.body.appendChild(f),f.click(),document.body.removeChild(f),URL.revokeObjectURL(p)}const Pe=()=>he("color-mode").value;var R={d:(e,t)=>{for(var n in t)R.o(t,n)&&!R.o(e,n)&&Object.defineProperty(e,n,{enumerable:!0,get:t[n]})},o:(e,t)=>Object.prototype.hasOwnProperty.call(e,t)},ie={};function q(e,t){(t==null||t>e.length)&&(t=e.length);for(var n=0,r=new Array(t);n<t;n++)r[n]=e[n];return r}function se(e,t){if(e){if(typeof e=="string")return q(e,t);var n=Object.prototype.toString.call(e).slice(8,-1);return n==="Object"&&e.constructor&&(n=e.constructor.name),n==="Map"||n==="Set"?Array.from(e):n==="Arguments"||/^(?:Ui|I)nt(?:8|16|32)(?:Clamped)?Array$/.test(n)?q(e,t):void 0}}function K(e){return function(t){if(Array.isArray(t))return q(t)}(e)||function(t){if(typeof Symbol<"u"&&t[Symbol.iterator]!=null||t["@@iterator"]!=null)return Array.from(t)}(e)||se(e)||function(){throw new TypeError(`Invalid attempt to spread non-iterable instance.
In order to be iterable, non-array objects must have a [Symbol.iterator]() method.`)}()}function H(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}R.d(ie,{Z:()=>De});const a=(te={computed:()=>fe,createTextVNode:()=>D,createVNode:()=>T,defineComponent:()=>le,reactive:()=>me,ref:()=>ye,watch:()=>ce,watchEffect:()=>ge},Y={},R.d(Y,te),Y),Te=(0,a.defineComponent)({props:{data:{required:!0,type:String},onClick:Function},render:function(){var e=this.data,t=this.onClick;return(0,a.createVNode)("span",{class:"vjs-tree-brackets",onClick:t},[e])}}),Be=(0,a.defineComponent)({emits:["change","update:modelValue"],props:{checked:{type:Boolean,default:!1},isMultiple:Boolean,onChange:Function},setup:function(e,t){var n=t.emit;return{uiType:(0,a.computed)(function(){return e.isMultiple?"checkbox":"radio"}),model:(0,a.computed)({get:function(){return e.checked},set:function(r){return n("update:modelValue",r)}})}},render:function(){var e=this.uiType,t=this.model,n=this.$emit;return(0,a.createVNode)("label",{class:["vjs-check-controller",t?"is-checked":""],onClick:function(r){return r.stopPropagation()}},[(0,a.createVNode)("span",{class:"vjs-check-controller-inner is-".concat(e)},null),(0,a.createVNode)("input",{checked:t,class:"vjs-check-controller-original is-".concat(e),type:e,onChange:function(){return n("change",t)}},null)])}}),Le=(0,a.defineComponent)({props:{nodeType:{required:!0,type:String},onClick:Function},render:function(){var e=this.nodeType,t=this.onClick,n=e==="objectStart"||e==="arrayStart";return n||e==="objectCollapsed"||e==="arrayCollapsed"?(0,a.createVNode)("span",{class:"vjs-carets vjs-carets-".concat(n?"open":"close"),onClick:t},[(0,a.createVNode)("svg",{viewBox:"0 0 1024 1024",focusable:"false","data-icon":"caret-down",width:"1em",height:"1em",fill:"currentColor","aria-hidden":"true"},[(0,a.createVNode)("path",{d:"M840.4 300H183.6c-19.7 0-30.7 20.8-18.5 35l328.4 380.8c9.4 10.9 27.5 10.9 37 0L858.9 335c12.2-14.2 1.2-35-18.5-35z"},null)])]):null}});var te,Y;function Q(e){return Q=typeof Symbol=="function"&&typeof Symbol.iterator=="symbol"?function(t){return typeof t}:function(t){return t&&typeof Symbol=="function"&&t.constructor===Symbol&&t!==Symbol.prototype?"symbol":typeof t},Q(e)}function ue(e){return Object.prototype.toString.call(e).slice(8,-1).toLowerCase()}function B(e){var t=arguments.length>1&&arguments[1]!==void 0?arguments[1]:"root",n=arguments.length>2&&arguments[2]!==void 0?arguments[2]:0,r=arguments.length>3?arguments[3]:void 0,i=r||{},p=i.key,f=i.index,d=i.type,b=d===void 0?"content":d,g=i.showComma,N=g!==void 0&&g,w=i.length,j=w===void 0?1:w,S=ue(e);if(S==="array"){var V=ne(e.map(function(m,s,o){return B(m,"".concat(t,"[").concat(s,"]"),n+1,{index:s,showComma:s!==o.length-1,length:j,type:b})}));return[B("[",t,n,{showComma:!1,key:p,length:e.length,type:"arrayStart"})[0]].concat(V,B("]",t,n,{showComma:N,length:e.length,type:"arrayEnd"})[0])}if(S==="object"){var _=Object.keys(e),P=ne(_.map(function(m,s,o){return B(e[m],/^[a-zA-Z_]\w*$/.test(m)?"".concat(t,".").concat(m):"".concat(t,'["').concat(m,'"]'),n+1,{key:m,showComma:s!==o.length-1,length:j,type:b})}));return[B("{",t,n,{showComma:!1,key:p,index:f,length:_.length,type:"objectStart"})[0]].concat(P,B("}",t,n,{showComma:N,length:_.length,type:"objectEnd"})[0])}return[{content:e,level:n,key:p,index:f,path:t,showComma:N,length:j,type:b}]}function ne(e){if(typeof Array.prototype.flat=="function")return e.flat();for(var t=K(e),n=[];t.length;){var r=t.shift();Array.isArray(r)?t.unshift.apply(t,K(r)):n.push(r)}return n}function $(e){var t=arguments.length>1&&arguments[1]!==void 0?arguments[1]:new WeakMap;if(e==null)return e;if(e instanceof Date)return new Date(e);if(e instanceof RegExp)return new RegExp(e);if(Q(e)!=="object")return e;if(t.get(e))return t.get(e);if(Array.isArray(e)){var n=e.map(function(p){return $(p,t)});return t.set(e,n),n}var r={};for(var i in e)r[i]=$(e[i],t);return t.set(e,r),r}function oe(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter(function(i){return Object.getOwnPropertyDescriptor(e,i).enumerable})),n.push.apply(n,r)}return n}function re(e){for(var t=1;t<arguments.length;t++){var n=arguments[t]!=null?arguments[t]:{};t%2?oe(Object(n),!0).forEach(function(r){H(e,r,n[r])}):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):oe(Object(n)).forEach(function(r){Object.defineProperty(e,r,Object.getOwnPropertyDescriptor(n,r))})}return e}var de={showLength:{type:Boolean,default:!1},showDoubleQuotes:{type:Boolean,default:!0},renderNodeKey:Function,renderNodeValue:Function,selectableType:String,showSelectController:{type:Boolean,default:!1},showLine:{type:Boolean,default:!0},showLineNumber:{type:Boolean,default:!1},selectOnClickNode:{type:Boolean,default:!0},nodeSelectable:{type:Function,default:function(){return!0}},highlightSelectedNode:{type:Boolean,default:!0},showIcon:{type:Boolean,default:!1},theme:{type:String,default:"light"},showKeyValueSpace:{type:Boolean,default:!0},editable:{type:Boolean,default:!1},editableTrigger:{type:String,default:"click"},onNodeClick:{type:Function},onBracketsClick:{type:Function},onIconClick:{type:Function},onValueChange:{type:Function}};const Ae=(0,a.defineComponent)({name:"TreeNode",props:re(re({},de),{},{node:{type:Object,required:!0},collapsed:Boolean,checked:Boolean,style:Object,onSelectedChange:{type:Function}}),emits:["nodeClick","bracketsClick","iconClick","selectedChange","valueChange"],setup:function(e,t){var n=t.emit,r=(0,a.computed)(function(){return ue(e.node.content)}),i=(0,a.computed)(function(){return"vjs-value vjs-value-".concat(r.value)}),p=(0,a.computed)(function(){return e.showDoubleQuotes?'"'.concat(e.node.key,'"'):e.node.key}),f=(0,a.computed)(function(){return e.selectableType==="multiple"}),d=(0,a.computed)(function(){return e.selectableType==="single"}),b=(0,a.computed)(function(){return e.nodeSelectable(e.node)&&(f.value||d.value)}),g=(0,a.reactive)({editing:!1}),N=function(s){var o,l,c=(l=(o=s.target)===null||o===void 0?void 0:o.value)==="null"?null:l==="undefined"?void 0:l==="true"||l!=="false"&&(l[0]+l[l.length-1]==='""'||l[0]+l[l.length-1]==="''"?l.slice(1,-1):typeof Number(l)=="number"&&!isNaN(Number(l))||l==="NaN"?Number(l):l);n("valueChange",c,e.node.path)},w=(0,a.computed)(function(){var s,o=(s=e.node)===null||s===void 0?void 0:s.content;return o===null?o="null":o===void 0&&(o="undefined"),r.value==="string"?'"'.concat(o,'"'):o+""}),j=function(){var s=e.renderNodeValue;return s?s({node:e.node,defaultValue:w.value}):w.value},S=function(){n("bracketsClick",!e.collapsed,e.node)},V=function(){n("iconClick",!e.collapsed,e.node)},_=function(){n("selectedChange",e.node)},P=function(){n("nodeClick",e.node),b.value&&e.selectOnClickNode&&n("selectedChange",e.node)},m=function(s){if(e.editable&&!g.editing){g.editing=!0;var o=function l(c){var u;c.target!==s.target&&((u=c.target)===null||u===void 0?void 0:u.parentElement)!==s.target&&(g.editing=!1,document.removeEventListener("click",l))};document.removeEventListener("click",o),document.addEventListener("click",o)}};return function(){var s,o=e.node;return(0,a.createVNode)("div",{class:{"vjs-tree-node":!0,"has-selector":e.showSelectController,"has-carets":e.showIcon,"is-highlight":e.highlightSelectedNode&&e.checked,dark:e.theme==="dark"},onClick:P,style:e.style},[e.showLineNumber&&(0,a.createVNode)("span",{class:"vjs-node-index"},[o.id+1]),e.showSelectController&&b.value&&o.type!=="objectEnd"&&o.type!=="arrayEnd"&&(0,a.createVNode)(Be,{isMultiple:f.value,checked:e.checked,onChange:_},null),(0,a.createVNode)("div",{class:"vjs-indent"},[Array.from(Array(o.level)).map(function(l,c){return(0,a.createVNode)("div",{key:c,class:{"vjs-indent-unit":!0,"has-line":e.showLine}},null)}),e.showIcon&&(0,a.createVNode)(Le,{nodeType:o.type,onClick:V},null)]),o.key&&(0,a.createVNode)("span",{class:"vjs-key"},[(s=e.renderNodeKey,s?s({node:e.node,defaultKey:p.value||""}):p.value),(0,a.createVNode)("span",{class:"vjs-colon"},[":".concat(e.showKeyValueSpace?" ":"")])]),(0,a.createVNode)("span",null,[o.type!=="content"&&o.content?(0,a.createVNode)(Te,{data:o.content.toString(),onClick:S},null):(0,a.createVNode)("span",{class:i.value,onClick:!e.editable||e.editableTrigger&&e.editableTrigger!=="click"?void 0:m,onDblclick:e.editable&&e.editableTrigger==="dblclick"?m:void 0},[e.editable&&g.editing?(0,a.createVNode)("input",{value:w.value,onChange:N,style:{padding:"3px 8px",border:"1px solid #eee",boxShadow:"none",boxSizing:"border-box",borderRadius:5,fontFamily:"inherit"}},null):j()]),o.showComma&&(0,a.createVNode)("span",null,[","]),e.showLength&&e.collapsed&&(0,a.createVNode)("span",{class:"vjs-comment"},[(0,a.createTextVNode)(" // "),o.length,(0,a.createTextVNode)(" items ")])])])}}});function ae(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter(function(i){return Object.getOwnPropertyDescriptor(e,i).enumerable})),n.push.apply(n,r)}return n}function k(e){for(var t=1;t<arguments.length;t++){var n=arguments[t]!=null?arguments[t]:{};t%2?ae(Object(n),!0).forEach(function(r){H(e,r,n[r])}):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):ae(Object(n)).forEach(function(r){Object.defineProperty(e,r,Object.getOwnPropertyDescriptor(n,r))})}return e}const De=(0,a.defineComponent)({name:"Tree",props:k(k({},de),{},{data:{type:[String,Number,Boolean,Array,Object],default:null},collapsedNodeLength:{type:Number,default:1/0},deep:{type:Number,default:1/0},pathCollapsible:{type:Function,default:function(){return!1}},rootPath:{type:String,default:"root"},virtual:{type:Boolean,default:!1},height:{type:Number,default:400},itemHeight:{type:Number,default:20},selectedValue:{type:[String,Array],default:function(){return""}},collapsedOnClickBrackets:{type:Boolean,default:!0},style:Object,onSelectedChange:{type:Function},theme:{type:String,default:"light"}}),slots:["renderNodeKey","renderNodeValue"],emits:["nodeClick","bracketsClick","iconClick","selectedChange","update:selectedValue","update:data"],setup:function(e,t){var n=t.emit,r=t.slots,i=(0,a.ref)(),p=(0,a.computed)(function(){return B(e.data,e.rootPath)}),f=function(o,l){return p.value.reduce(function(c,u){var h,y=u.level>=o||u.length>=l,C=(h=e.pathCollapsible)===null||h===void 0?void 0:h.call(e,u);return u.type!=="objectStart"&&u.type!=="arrayStart"||!y&&!C?c:k(k({},c),{},H({},u.path,1))},{})},d=(0,a.reactive)({translateY:0,visibleData:null,hiddenPaths:f(e.deep,e.collapsedNodeLength)}),b=(0,a.computed)(function(){for(var o=null,l=[],c=p.value.length,u=0;u<c;u++){var h=k(k({},p.value[u]),{},{id:u}),y=d.hiddenPaths[h.path];if(o&&o.path===h.path){var C=o.type==="objectStart",I=k(k(k({},h),o),{},{showComma:h.showComma,content:C?"{...}":"[...]",type:C?"objectCollapsed":"arrayCollapsed"});o=null,l.push(I)}else{if(y&&!o){o=h;continue}if(o)continue;l.push(h)}}return l}),g=(0,a.computed)(function(){var o=e.selectedValue;return o&&e.selectableType==="multiple"&&Array.isArray(o)?o:[o]}),N=(0,a.computed)(function(){return!e.selectableType||e.selectOnClickNode||e.showSelectController?"":"When selectableType is not null, selectOnClickNode and showSelectController cannot be false at the same time, because this will cause the selection to fail."}),w=function(){var o=b.value;if(e.virtual){var l,c=e.height/e.itemHeight,u=((l=i.value)===null||l===void 0?void 0:l.scrollTop)||0,h=Math.floor(u/e.itemHeight),y=h<0?0:h+c>o.length?o.length-c:h;y<0&&(y=0);var C=y+c;d.translateY=y*e.itemHeight,d.visibleData=o.filter(function(I,E){return E>=y&&E<C})}else d.visibleData=o},j=function(){w()},S=function(o){var l,c,u=o.path,h=e.selectableType;if(h==="multiple"){var y=g.value.findIndex(function(O){return O===u}),C=K(g.value);y!==-1?C.splice(y,1):C.push(u),n("update:selectedValue",C),n("selectedChange",C,K(g.value))}else if(h==="single"&&g.value[0]!==u){var I=(l=g.value,c=1,function(O){if(Array.isArray(O))return O}(l)||function(O,J){var L=O==null?null:typeof Symbol<"u"&&O[Symbol.iterator]||O["@@iterator"];if(L!=null){var Z,G,U=[],z=!0,W=!1;try{for(L=L.call(O);!(z=(Z=L.next()).done)&&(U.push(Z.value),!J||U.length!==J);z=!0);}catch(pe){W=!0,G=pe}finally{try{z||L.return==null||L.return()}finally{if(W)throw G}}return U}}(l,c)||se(l,c)||function(){throw new TypeError(`Invalid attempt to destructure non-iterable instance.
In order to be iterable, non-array objects must have a [Symbol.iterator]() method.`)}())[0],E=u;n("update:selectedValue",E),n("selectedChange",E,I)}},V=function(o){n("nodeClick",o)},_=function(o,l){if(o)d.hiddenPaths=k(k({},d.hiddenPaths),{},H({},l,1));else{var c=k({},d.hiddenPaths);delete c[l],d.hiddenPaths=c}},P=function(o,l){e.collapsedOnClickBrackets&&_(o,l.path),n("bracketsClick",o,l)},m=function(o,l){_(o,l.path),n("iconClick",o,l)},s=function(o,l){var c=$(e.data),u=e.rootPath;new Function("data","val","data".concat(l.slice(u.length),"=val"))(c,o),n("update:data",c)};return(0,a.watchEffect)(function(){N.value&&function(o){throw new Error("[VueJSONPretty] ".concat(o))}(N.value)}),(0,a.watchEffect)(function(){b.value&&w()}),(0,a.watch)(function(){return e.deep},function(o){o&&(d.hiddenPaths=f(o,e.collapsedNodeLength))}),(0,a.watch)(function(){return e.collapsedNodeLength},function(o){o&&(d.hiddenPaths=f(e.deep,o))}),function(){var o,l,c=(o=e.renderNodeKey)!==null&&o!==void 0?o:r.renderNodeKey,u=(l=e.renderNodeValue)!==null&&l!==void 0?l:r.renderNodeValue,h=d.visibleData&&d.visibleData.map(function(y){return(0,a.createVNode)(Ae,{key:y.id,node:y,collapsed:!!d.hiddenPaths[y.path],theme:e.theme,showDoubleQuotes:e.showDoubleQuotes,showLength:e.showLength,checked:g.value.includes(y.path),selectableType:e.selectableType,showLine:e.showLine,showLineNumber:e.showLineNumber,showSelectController:e.showSelectController,selectOnClickNode:e.selectOnClickNode,nodeSelectable:e.nodeSelectable,highlightSelectedNode:e.highlightSelectedNode,editable:e.editable,editableTrigger:e.editableTrigger,showIcon:e.showIcon,showKeyValueSpace:e.showKeyValueSpace,renderNodeKey:c,renderNodeValue:u,onNodeClick:V,onBracketsClick:P,onIconClick:m,onSelectedChange:S,onValueChange:s,style:e.itemHeight&&e.itemHeight!==20?{lineHeight:"".concat(e.itemHeight,"px")}:{}},null)});return(0,a.createVNode)("div",{ref:i,class:{"vjs-tree":!0,"is-virtual":e.virtual,dark:e.theme==="dark"},onScroll:e.virtual?j:void 0,style:e.showLineNumber?k({paddingLeft:"".concat(12*Number(p.value.length.toString().length),"px")},e.style):e.style},[e.virtual?(0,a.createVNode)("div",{class:"vjs-tree-list",style:{height:"".concat(e.height,"px")}},[(0,a.createVNode)("div",{class:"vjs-tree-list-holder",style:{height:"".concat(b.value.length*e.itemHeight,"px")}},[(0,a.createVNode)("div",{class:"vjs-tree-list-holder-inner",style:{transform:"translateY(".concat(d.translateY,"px)")}},[h])])]):h])}}});var Ee=ie.Z;const Fe={class:"flex flex-col flex-1 gap-y-4 max-h-screen overflow-auto p-4"},Ie={class:"flex items-center justify-between"},Me=x("p",{class:"text-base font-semibold leading-6 text-gray-900 dark:text-white"}," Document Sample ",-1),Re={key:0},Ke={key:1,class:"flex items-center gap-x-2 text-lg"},He={key:3,class:"text-red-600 flex flex-col gap-y-2 items-center"},Ue={class:"font-semibold w-full"},ze=x("p",{class:"w-full"},"There was an error retrieving the document:",-1),Ye=x("p",null,"If the problem persists, please verify that the data mapping might be incorrect or that the server is responding as expected.",-1),qe=x("p",null,"You may need to check the response format or any related configuration.",-1),Qe={class:"font-semibold text-lg w-full"},$e={class:"flex items-center justify-between"},Je={class:"flex items-center gap-x-4 cursor-pointer"},Xe=le({__name:"SamplePanel",props:{form:{required:!0},formModifiers:{}},emits:["update:form"],async setup(e){let t,n;const r=je(),i=ve(e,"form"),{status:p,error:f,data:d,refresh:b}=([t,n]=be(()=>Oe("/processors/test",{method:"post",body:i,watch:!1,responseType:"json",immediate:!1})),t=await t,n(),t),g=xe(async()=>{await b()},4e3);ce(i.value,async w=>{r.value.liveReload&&await g()});const N=()=>{d.value&&Ve(d.value,i.value.index+"_sample.json")};return(w,j)=>{const S=Ce,V=Ne,_=Se,P=_e;return F(),M("div",Fe,[T(P,{class:"flex flex-col flex-1 overflow-x-hidden overflow-y-auto",ui:{divide:"",body:{base:"flex-1"},header:{base:"sticky z-20 top-0 shadow-md bg-white dark:bg-gray-900 border-b border-gray-200 dark:border-gray-800"},footer:{base:"sticky z-20 bottom-0 bg-white dark:bg-gray-900 border-t border-gray-200 dark:border-gray-800"}}},{header:A(()=>[x("div",Ie,[Me,T(S,{color:"sky",loading:v(p)==="pending",onClick:X(v(b),["prevent"])},{default:A(()=>[D(" Fetch Sample ")]),_:1},8,["loading","onClick"])])]),footer:A(()=>[x("div",$e,[x("label",Je,[T(_,{"on-icon":"i-heroicons-check-20-solid","off-icon":"i-heroicons-x-mark-20-solid",modelValue:v(r).liveReload,"onUpdate:modelValue":j[0]||(j[0]=m=>v(r).liveReload=m)},null,8,["modelValue"]),D(" Live reload ")]),T(S,{color:"gray",disabled:!v(d),onClick:X(N,["prevent","stop"])},{default:A(()=>[D(" Download ")]),_:1},8,["disabled"])])]),default:A(()=>{var m,s,o;return[v(p)==="idle"?(F(),M("div",Re,' Please click the "Fetch Sample" button to generate the document. ')):v(p)==="pending"?(F(),M("div",Ke,[D(" Generating document "),T(V,{name:"i-heroicons-arrow-path",class:"animate-spin"})])):v(p)=="success"?(F(),we(v(Ee),{key:2,"show-icon":"","show-length":"","show-line-number":"",data:v(d),theme:("useColorMode"in w?w.useColorMode:v(Pe))().value},null,8,["data","theme"])):v(p)==="error"?(F(),M("div",He,[x("p",Ue,ee((s=(m=v(f))==null?void 0:m.data)==null?void 0:s.error),1),ze,Ye,qe,x("p",Qe,ee((o=v(f))==null?void 0:o.message),1),T(S,{color:"red",onClick:v(b)},{default:A(()=>[D("Retry")]),_:1},8,["onClick"])])):ke("",!0)]}),_:1})])}}});export{Xe as default};
