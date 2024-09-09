import{ap as at,aq as ne,r as g,ar as nt,U as Te,as as Pt,H as k,M as Z,at as Xe,au as Ct,h as Ot,av as Ue,aw as At,s as m,ax as Y,ay as Rt,am as rt,az as _t,K as R,aA as Nt,x as S,A as q,J as lt,v as _,g as A,C as x,G as O,O as L,F as Mt,y as I,aB as kt,I as Ve,E as ot,B as C,N as _e,L as Bt,aC as Ht,aD as J,z as It,a6 as le,T as xt,_ as it,Q as st,R as Se,S as ut,V as jt,W as Ut,o as re,a5 as Ne,w as ae,c as Ce,X as be,m as we,Y as Ee,a9 as Vt,a2 as qt,a4 as me,b as Oe,a as Wt,aE as zt}from"./CcBA2A6A.js";const Kt=e=>e==="defer"||e===!1;function Gt(...e){var E;const t=typeof e[e.length-1]=="string"?e.pop():void 0;typeof e[0]!="string"&&e.unshift(t);let[a,r,l={}]=e;if(typeof a!="string")throw new TypeError("[nuxt] [asyncData] key must be a string.");if(typeof r!="function")throw new TypeError("[nuxt] [asyncData] handler must be a function.");const n=at(),o=r,i=()=>ne.value,s=()=>n.isHydrating?n.payload.data[a]:n.static.data[a];l.server=l.server??!0,l.default=l.default??i,l.getCachedData=l.getCachedData??s,l.lazy=l.lazy??!1,l.immediate=l.immediate??!0,l.deep=l.deep??ne.deep,l.dedupe=l.dedupe??"cancel";const u=l.getCachedData(a,n),d=u!=null;if(!n._asyncData[a]||!l.immediate){(E=n.payload._errors)[a]??(E[a]=ne.errorValue);const p=l.deep?g:nt;n._asyncData[a]={data:p(d?u:l.default()),pending:g(!d),error:Te(n.payload._errors,a),status:g("idle"),_default:l.default}}const c={...n._asyncData[a]};delete c._default,c.refresh=c.execute=(p={})=>{if(n._asyncDataPromises[a]){if(Kt(p.dedupe??l.dedupe))return n._asyncDataPromises[a];n._asyncDataPromises[a].cancelled=!0}if(p._initial||n.isHydrating&&p._initial!==!1){const b=p._initial?u:l.getCachedData(a,n);if(b!=null)return Promise.resolve(b)}c.pending.value=!0,c.status.value="pending";const $=new Promise((b,y)=>{try{b(o(n))}catch(N){y(N)}}).then(async b=>{if($.cancelled)return n._asyncDataPromises[a];let y=b;l.transform&&(y=await l.transform(b)),l.pick&&(y=Xt(y,l.pick)),n.payload.data[a]=y,c.data.value=y,c.error.value=ne.errorValue,c.status.value="success"}).catch(b=>{if($.cancelled)return n._asyncDataPromises[a];c.error.value=Ct(b),c.data.value=Ot(l.default()),c.status.value="error"}).finally(()=>{$.cancelled||(c.pending.value=!1,delete n._asyncDataPromises[a])});return n._asyncDataPromises[a]=$,n._asyncDataPromises[a]},c.clear=()=>Yt(n,a);const f=()=>c.refresh({_initial:!0}),h=l.server!==!1&&n.payload.serverRendered;{const p=Ue();if(p&&!p._nuxtOnBeforeMountCbs){p._nuxtOnBeforeMountCbs=[];const y=p._nuxtOnBeforeMountCbs;Pt(()=>{y.forEach(N=>{N()}),y.splice(0,y.length)}),k(()=>y.splice(0,y.length))}h&&n.isHydrating&&(c.error.value||u!=null)?(c.pending.value=!1,c.status.value=c.error.value?"error":"success"):p&&(n.payload.serverRendered&&n.isHydrating||l.lazy)&&l.immediate?p._nuxtOnBeforeMountCbs.push(f):l.immediate&&f();const $=At();if(l.watch){const y=Z(l.watch,()=>c.refresh());$&&Xe(y)}const b=n.hook("app:data:refresh",async y=>{(!y||y.includes(a))&&await c.refresh()});$&&Xe(b)}const v=Promise.resolve(n._asyncDataPromises[a]).then(()=>c);return Object.assign(v,c),v}function Yt(e,t){t in e.payload.data&&(e.payload.data[t]=void 0),t in e.payload._errors&&(e.payload._errors[t]=ne.errorValue),e._asyncData[t]&&(e._asyncData[t].data.value=void 0,e._asyncData[t].error.value=ne.errorValue,e._asyncData[t].pending.value=!1,e._asyncData[t].status.value="idle"),t in e._asyncDataPromises&&(e._asyncDataPromises[t]&&(e._asyncDataPromises[t].cancelled=!0),e._asyncDataPromises[t]=void 0)}function Xt(e,t){const a={};for(const r of t)a[r]=e[r];return a}function Jt(e,t,a){const[r={},l]=typeof t=="string"?[{},t]:[t,a],n=m(()=>Y(e)),o=r.key||Rt([l,typeof n.value=="string"?n.value:"",...Qt(r)]);if(!o||typeof o!="string")throw new TypeError("[nuxt] [useFetch] key must be a string: "+o);if(!e)throw new Error("[nuxt] [useFetch] request is missing.");const i=o===l?"$f"+o:o;if(!r.baseURL&&typeof n.value=="string"&&n.value[0]==="/"&&n.value[1]==="/")throw new Error('[nuxt] [useFetch] the request URL must not start with "//".');const{server:s,lazy:u,default:d,transform:c,pick:f,watch:h,immediate:v,getCachedData:E,deep:p,dedupe:$,...b}=r,y=rt({..._t,...b,cache:typeof r.cache=="boolean"?void 0:r.cache}),N={server:s,lazy:u,default:d,transform:c,pick:f,immediate:v,getCachedData:E,deep:p,dedupe:$,watch:h===!1?[]:[y,n,...h||[]]};let D;return Gt(i,()=>{var K;(K=D==null?void 0:D.abort)==null||K.call(D,"Request aborted as another request to the same endpoint was initiated."),D=typeof AbortController<"u"?new AbortController:{};const z=Y(r.timeout);let j;return z&&(j=setTimeout(()=>D.abort("Request aborted due to timeout."),z),D.signal.onabort=()=>clearTimeout(j)),(r.$fetch||globalThis.$fetch)(n.value,{signal:D.signal,...y}).finally(()=>{clearTimeout(j)})},N)}function Qt(e){var a;const t=[((a=Y(e.method))==null?void 0:a.toUpperCase())||"GET",Y(e.baseURL)];for(const r of[e.params||e.query]){const l=Y(r);if(!l)continue;const n={};for(const[o,i]of Object.entries(l))n[Y(o)]=Y(i);t.push(n)}return t}const Zt="data-n-ids",ea="-";function ta(e){var l,n,o,i,s,u;if(typeof e!="string")throw new TypeError("[nuxt] [useId] key must be a string.");e=`n${e.slice(1)}`;const t=at(),a=Ue();if(!a)throw new TypeError("[nuxt] `useId` must be called within a component setup function.");t._genId||(t._genId=0),a._nuxtIdIndex||(a._nuxtIdIndex={}),(l=a._nuxtIdIndex)[e]||(l[e]=0);const r=e+ea+a._nuxtIdIndex[e]++;if(t.payload.serverRendered&&t.isHydrating&&!R(Nt,!1)){const d=((n=a.vnode.el)==null?void 0:n.nodeType)===8&&((i=(o=a.vnode.el)==null?void 0:o.nextElementSibling)!=null&&i.getAttribute)?(s=a.vnode.el)==null?void 0:s.nextElementSibling:a.vnode.el,c=JSON.parse(((u=d==null?void 0:d.getAttribute)==null?void 0:u.call(d,Zt))||"{}");if(c[r])return c[r]}return e+"_"+t._genId++}const aa={base:"",background:"bg-white dark:bg-gray-900",divide:"divide-y divide-gray-200 dark:divide-gray-800",ring:"ring-1 ring-gray-200 dark:ring-gray-800",rounded:"rounded-lg",shadow:"shadow",body:{base:"",background:"",padding:"px-4 py-5 sm:p-6"},header:{base:"",background:"",padding:"px-4 py-5 sm:px-6"},footer:{base:"",background:"",padding:"px-4 py-4 sm:px-6"}},na={wrapper:"fixed inset-0 flex z-50",overlay:{base:"fixed inset-0 transition-opacity",background:"bg-gray-200/75 dark:bg-gray-800/75",transition:{enter:"ease-in-out duration-500",enterFrom:"opacity-0",enterTo:"opacity-100",leave:"ease-in-out duration-500",leaveFrom:"opacity-100",leaveTo:"opacity-0"}},base:"relative flex-1 flex flex-col w-full focus:outline-none",background:"bg-white dark:bg-gray-900",ring:"",rounded:"",padding:"",shadow:"shadow-xl",width:"w-screen max-w-md",height:"h-screen max-h-96",translate:{base:"translate-x-0",left:"-translate-x-full rtl:translate-x-full",right:"translate-x-full rtl:-translate-x-full",top:"-translate-y-full",bottom:"translate-y-full"},transition:{enter:"transform transition ease-in-out duration-300",leave:"transform transition ease-in-out duration-200"}};function qe(e){typeof queueMicrotask=="function"?queueMicrotask(e):Promise.resolve().then(e).catch(t=>setTimeout(()=>{throw t}))}function ce(){let e=[],t={addEventListener(a,r,l,n){return a.addEventListener(r,l,n),t.add(()=>a.removeEventListener(r,l,n))},requestAnimationFrame(...a){let r=requestAnimationFrame(...a);t.add(()=>cancelAnimationFrame(r))},nextFrame(...a){t.requestAnimationFrame(()=>{t.requestAnimationFrame(...a)})},setTimeout(...a){let r=setTimeout(...a);t.add(()=>clearTimeout(r))},microTask(...a){let r={current:!0};return qe(()=>{r.current&&a[0]()}),t.add(()=>{r.current=!1})},style(a,r,l){let n=a.style.getPropertyValue(r);return Object.assign(a.style,{[r]:l}),this.add(()=>{Object.assign(a.style,{[r]:n})})},group(a){let r=ce();return a(r),this.add(()=>r.dispose())},add(a){return e.push(a),()=>{let r=e.indexOf(a);if(r>=0)for(let l of e.splice(r,1))l()}},dispose(){for(let a of e.splice(0))a()}};return t}var ra=Object.defineProperty,la=(e,t,a)=>t in e?ra(e,t,{enumerable:!0,configurable:!0,writable:!0,value:a}):e[t]=a,Je=(e,t,a)=>(la(e,typeof t!="symbol"?t+"":t,a),a);let oa=class{constructor(){Je(this,"current",this.detect()),Je(this,"currentId",0)}set(t){this.current!==t&&(this.currentId=0,this.current=t)}reset(){this.set(this.detect())}nextId(){return++this.currentId}get isServer(){return this.current==="server"}get isClient(){return this.current==="client"}detect(){return typeof window>"u"||typeof document>"u"?"server":"client"}},fe=new oa;function ee(e){if(fe.isServer)return null;if(e instanceof Node)return e.ownerDocument;if(e!=null&&e.hasOwnProperty("value")){let t=S(e);if(t)return t.ownerDocument}return document}let Me=["[contentEditable=true]","[tabindex]","a[href]","area[href]","button:not([disabled])","iframe","input:not([disabled])","select:not([disabled])","textarea:not([disabled])"].map(e=>`${e}:not([tabindex='-1'])`).join(",");var U=(e=>(e[e.First=1]="First",e[e.Previous=2]="Previous",e[e.Next=4]="Next",e[e.Last=8]="Last",e[e.WrapAround=16]="WrapAround",e[e.NoScroll=32]="NoScroll",e))(U||{}),dt=(e=>(e[e.Error=0]="Error",e[e.Overflow=1]="Overflow",e[e.Success=2]="Success",e[e.Underflow=3]="Underflow",e))(dt||{}),ia=(e=>(e[e.Previous=-1]="Previous",e[e.Next=1]="Next",e))(ia||{});function ct(e=document.body){return e==null?[]:Array.from(e.querySelectorAll(Me)).sort((t,a)=>Math.sign((t.tabIndex||Number.MAX_SAFE_INTEGER)-(a.tabIndex||Number.MAX_SAFE_INTEGER)))}var ft=(e=>(e[e.Strict=0]="Strict",e[e.Loose=1]="Loose",e))(ft||{});function vt(e,t=0){var a;return e===((a=ee(e))==null?void 0:a.body)?!1:q(t,{0(){return e.matches(Me)},1(){let r=e;for(;r!==null;){if(r.matches(Me))return!0;r=r.parentElement}return!1}})}function nn(e){let t=ee(e);lt(()=>{t&&!vt(t.activeElement,0)&&V(e)})}var sa=(e=>(e[e.Keyboard=0]="Keyboard",e[e.Mouse=1]="Mouse",e))(sa||{});typeof window<"u"&&typeof document<"u"&&(document.addEventListener("keydown",e=>{e.metaKey||e.altKey||e.ctrlKey||(document.documentElement.dataset.headlessuiFocusVisible="")},!0),document.addEventListener("click",e=>{e.detail===1?delete document.documentElement.dataset.headlessuiFocusVisible:e.detail===0&&(document.documentElement.dataset.headlessuiFocusVisible="")},!0));function V(e){e==null||e.focus({preventScroll:!0})}let ua=["textarea","input"].join(",");function da(e){var t,a;return(a=(t=e==null?void 0:e.matches)==null?void 0:t.call(e,ua))!=null?a:!1}function ca(e,t=a=>a){return e.slice().sort((a,r)=>{let l=t(a),n=t(r);if(l===null||n===null)return 0;let o=l.compareDocumentPosition(n);return o&Node.DOCUMENT_POSITION_FOLLOWING?-1:o&Node.DOCUMENT_POSITION_PRECEDING?1:0})}function rn(e,t){return de(ct(),t,{relativeTo:e})}function de(e,t,{sorted:a=!0,relativeTo:r=null,skipElements:l=[]}={}){var n;let o=(n=Array.isArray(e)?e.length>0?e[0].ownerDocument:document:e==null?void 0:e.ownerDocument)!=null?n:document,i=Array.isArray(e)?a?ca(e):e:ct(e);l.length>0&&i.length>1&&(i=i.filter(v=>!l.includes(v))),r=r??o.activeElement;let s=(()=>{if(t&5)return 1;if(t&10)return-1;throw new Error("Missing Focus.First, Focus.Previous, Focus.Next or Focus.Last")})(),u=(()=>{if(t&1)return 0;if(t&2)return Math.max(0,i.indexOf(r))-1;if(t&4)return Math.max(0,i.indexOf(r))+1;if(t&8)return i.length-1;throw new Error("Missing Focus.First, Focus.Previous, Focus.Next or Focus.Last")})(),d=t&32?{preventScroll:!0}:{},c=0,f=i.length,h;do{if(c>=f||c+f<=0)return 0;let v=u+c;if(t&16)v=(v+f)%f;else{if(v<0)return 3;if(v>=f)return 1}h=i[v],h==null||h.focus(d),c+=s}while(h!==o.activeElement);return t&6&&da(h)&&h.select(),2}function pt(){return/iPhone/gi.test(window.navigator.platform)||/Mac/gi.test(window.navigator.platform)&&window.navigator.maxTouchPoints>0}function fa(){return/Android/gi.test(window.navigator.userAgent)}function va(){return pt()||fa()}function he(e,t,a){fe.isServer||_(r=>{document.addEventListener(e,t,a),r(()=>document.removeEventListener(e,t,a))})}function mt(e,t,a){fe.isServer||_(r=>{window.addEventListener(e,t,a),r(()=>window.removeEventListener(e,t,a))})}function pa(e,t,a=m(()=>!0)){function r(n,o){if(!a.value||n.defaultPrevented)return;let i=o(n);if(i===null||!i.getRootNode().contains(i))return;let s=function u(d){return typeof d=="function"?u(d()):Array.isArray(d)||d instanceof Set?d:[d]}(e);for(let u of s){if(u===null)continue;let d=u instanceof HTMLElement?u:S(u);if(d!=null&&d.contains(i)||n.composed&&n.composedPath().includes(d))return}return!vt(i,ft.Loose)&&i.tabIndex!==-1&&n.preventDefault(),t(n,i)}let l=g(null);he("pointerdown",n=>{var o,i;a.value&&(l.value=((i=(o=n.composedPath)==null?void 0:o.call(n))==null?void 0:i[0])||n.target)},!0),he("mousedown",n=>{var o,i;a.value&&(l.value=((i=(o=n.composedPath)==null?void 0:o.call(n))==null?void 0:i[0])||n.target)},!0),he("click",n=>{va()||l.value&&(r(n,()=>l.value),l.value=null)},!0),he("touchend",n=>r(n,()=>n.target instanceof HTMLElement?n.target:null),!0),mt("blur",n=>r(n,()=>window.document.activeElement instanceof HTMLIFrameElement?window.document.activeElement:null),!0)}var $e=(e=>(e[e.None=1]="None",e[e.Focusable=2]="Focusable",e[e.Hidden=4]="Hidden",e))($e||{});let ke=A({name:"Hidden",props:{as:{type:[Object,String],default:"div"},features:{type:Number,default:1}},setup(e,{slots:t,attrs:a}){return()=>{var r;let{features:l,...n}=e,o={"aria-hidden":(l&2)===2?!0:(r=n["aria-hidden"])!=null?r:void 0,hidden:(l&4)===4?!0:void 0,style:{position:"fixed",top:1,left:1,width:1,height:0,padding:0,margin:-1,overflow:"hidden",clip:"rect(0, 0, 0, 0)",whiteSpace:"nowrap",borderWidth:"0",...(l&4)===4&&(l&2)!==2&&{display:"none"}}};return x({ourProps:o,theirProps:n,slot:{},attrs:a,slots:t,name:"Hidden"})}}});function ma(e){function t(){document.readyState!=="loading"&&(e(),document.removeEventListener("DOMContentLoaded",t))}typeof window<"u"&&typeof document<"u"&&(document.addEventListener("DOMContentLoaded",t),t())}let X=[];ma(()=>{function e(t){t.target instanceof HTMLElement&&t.target!==document.body&&X[0]!==t.target&&(X.unshift(t.target),X=X.filter(a=>a!=null&&a.isConnected),X.splice(10))}window.addEventListener("click",e,{capture:!0}),window.addEventListener("mousedown",e,{capture:!0}),window.addEventListener("focus",e,{capture:!0}),document.body.addEventListener("click",e,{capture:!0}),document.body.addEventListener("mousedown",e,{capture:!0}),document.body.addEventListener("focus",e,{capture:!0})});function ht(e,t,a,r){fe.isServer||_(l=>{e=e??window,e.addEventListener(t,a,r),l(()=>e.removeEventListener(t,a,r))})}var ue=(e=>(e[e.Forwards=0]="Forwards",e[e.Backwards=1]="Backwards",e))(ue||{});function ha(){let e=g(0);return mt("keydown",t=>{t.key==="Tab"&&(e.value=t.shiftKey?1:0)}),e}function gt(e){if(!e)return new Set;if(typeof e=="function")return new Set(e());let t=new Set;for(let a of e.value){let r=S(a);r instanceof HTMLElement&&t.add(r)}return t}var yt=(e=>(e[e.None=1]="None",e[e.InitialFocus=2]="InitialFocus",e[e.TabLock=4]="TabLock",e[e.FocusLock=8]="FocusLock",e[e.RestoreFocus=16]="RestoreFocus",e[e.All=30]="All",e))(yt||{});let ie=Object.assign(A({name:"FocusTrap",props:{as:{type:[Object,String],default:"div"},initialFocus:{type:Object,default:null},features:{type:Number,default:30},containers:{type:[Object,Function],default:g(new Set)}},inheritAttrs:!1,setup(e,{attrs:t,slots:a,expose:r}){let l=g(null);r({el:l,$el:l});let n=m(()=>ee(l)),o=g(!1);O(()=>o.value=!0),k(()=>o.value=!1),ya({ownerDocument:n},m(()=>o.value&&!!(e.features&16)));let i=ba({ownerDocument:n,container:l,initialFocus:m(()=>e.initialFocus)},m(()=>o.value&&!!(e.features&2)));wa({ownerDocument:n,container:l,containers:e.containers,previousActiveElement:i},m(()=>o.value&&!!(e.features&8)));let s=ha();function u(h){let v=S(l);v&&(E=>E())(()=>{q(s.value,{[ue.Forwards]:()=>{de(v,U.First,{skipElements:[h.relatedTarget]})},[ue.Backwards]:()=>{de(v,U.Last,{skipElements:[h.relatedTarget]})}})})}let d=g(!1);function c(h){h.key==="Tab"&&(d.value=!0,requestAnimationFrame(()=>{d.value=!1}))}function f(h){if(!o.value)return;let v=gt(e.containers);S(l)instanceof HTMLElement&&v.add(S(l));let E=h.relatedTarget;E instanceof HTMLElement&&E.dataset.headlessuiFocusGuard!=="true"&&(bt(v,E)||(d.value?de(S(l),q(s.value,{[ue.Forwards]:()=>U.Next,[ue.Backwards]:()=>U.Previous})|U.WrapAround,{relativeTo:h.target}):h.target instanceof HTMLElement&&V(h.target)))}return()=>{let h={},v={ref:l,onKeydown:c,onFocusout:f},{features:E,initialFocus:p,containers:$,...b}=e;return L(Mt,[!!(E&4)&&L(ke,{as:"button",type:"button","data-headlessui-focus-guard":!0,onFocus:u,features:$e.Focusable}),x({ourProps:v,theirProps:{...t,...b},slot:h,attrs:t,slots:a,name:"FocusTrap"}),!!(E&4)&&L(ke,{as:"button",type:"button","data-headlessui-focus-guard":!0,onFocus:u,features:$e.Focusable})])}}}),{features:yt});function ga(e){let t=g(X.slice());return Z([e],([a],[r])=>{r===!0&&a===!1?qe(()=>{t.value.splice(0)}):r===!1&&a===!0&&(t.value=X.slice())},{flush:"post"}),()=>{var a;return(a=t.value.find(r=>r!=null&&r.isConnected))!=null?a:null}}function ya({ownerDocument:e},t){let a=ga(t);O(()=>{_(()=>{var r,l;t.value||((r=e.value)==null?void 0:r.activeElement)===((l=e.value)==null?void 0:l.body)&&V(a())},{flush:"post"})}),k(()=>{t.value&&V(a())})}function ba({ownerDocument:e,container:t,initialFocus:a},r){let l=g(null),n=g(!1);return O(()=>n.value=!0),k(()=>n.value=!1),O(()=>{Z([t,a,r],(o,i)=>{if(o.every((u,d)=>(i==null?void 0:i[d])===u)||!r.value)return;let s=S(t);s&&qe(()=>{var u,d;if(!n.value)return;let c=S(a),f=(u=e.value)==null?void 0:u.activeElement;if(c){if(c===f){l.value=f;return}}else if(s.contains(f)){l.value=f;return}c?V(c):de(s,U.First|U.NoScroll)===dt.Error&&console.warn("There are no focusable elements inside the <FocusTrap />"),l.value=(d=e.value)==null?void 0:d.activeElement})},{immediate:!0,flush:"post"})}),l}function wa({ownerDocument:e,container:t,containers:a,previousActiveElement:r},l){var n;ht((n=e.value)==null?void 0:n.defaultView,"focus",o=>{if(!l.value)return;let i=gt(a);S(t)instanceof HTMLElement&&i.add(S(t));let s=r.value;if(!s)return;let u=o.target;u&&u instanceof HTMLElement?bt(i,u)?(r.value=u,V(u)):(o.preventDefault(),o.stopPropagation(),V(s)):V(r.value)},!0)}function bt(e,t){for(let a of e)if(a.contains(t))return!0;return!1}function Ea(e){let t=nt(e.getSnapshot());return k(e.subscribe(()=>{t.value=e.getSnapshot()})),t}function Ta(e,t){let a=e(),r=new Set;return{getSnapshot(){return a},subscribe(l){return r.add(l),()=>r.delete(l)},dispatch(l,...n){let o=t[l].call(a,...n);o&&(a=o,r.forEach(i=>i()))}}}function Sa(){let e;return{before({doc:t}){var a;let r=t.documentElement;e=((a=t.defaultView)!=null?a:window).innerWidth-r.clientWidth},after({doc:t,d:a}){let r=t.documentElement,l=r.clientWidth-r.offsetWidth,n=e-l;a.style(r,"paddingRight",`${n}px`)}}}function $a(){return pt()?{before({doc:e,d:t,meta:a}){function r(l){return a.containers.flatMap(n=>n()).some(n=>n.contains(l))}t.microTask(()=>{var l;if(window.getComputedStyle(e.documentElement).scrollBehavior!=="auto"){let i=ce();i.style(e.documentElement,"scrollBehavior","auto"),t.add(()=>t.microTask(()=>i.dispose()))}let n=(l=window.scrollY)!=null?l:window.pageYOffset,o=null;t.addEventListener(e,"click",i=>{if(i.target instanceof HTMLElement)try{let s=i.target.closest("a");if(!s)return;let{hash:u}=new URL(s.href),d=e.querySelector(u);d&&!r(d)&&(o=d)}catch{}},!0),t.addEventListener(e,"touchstart",i=>{if(i.target instanceof HTMLElement)if(r(i.target)){let s=i.target;for(;s.parentElement&&r(s.parentElement);)s=s.parentElement;t.style(s,"overscrollBehavior","contain")}else t.style(i.target,"touchAction","none")}),t.addEventListener(e,"touchmove",i=>{if(i.target instanceof HTMLElement){if(i.target.tagName==="INPUT")return;if(r(i.target)){let s=i.target;for(;s.parentElement&&s.dataset.headlessuiPortal!==""&&!(s.scrollHeight>s.clientHeight||s.scrollWidth>s.clientWidth);)s=s.parentElement;s.dataset.headlessuiPortal===""&&i.preventDefault()}else i.preventDefault()}},{passive:!1}),t.add(()=>{var i;let s=(i=window.scrollY)!=null?i:window.pageYOffset;n!==s&&window.scrollTo(0,n),o&&o.isConnected&&(o.scrollIntoView({block:"nearest"}),o=null)})})}}:{}}function Da(){return{before({doc:e,d:t}){t.style(e.documentElement,"overflow","hidden")}}}function La(e){let t={};for(let a of e)Object.assign(t,a(t));return t}let Q=Ta(()=>new Map,{PUSH(e,t){var a;let r=(a=this.get(e))!=null?a:{doc:e,count:0,d:ce(),meta:new Set};return r.count++,r.meta.add(t),this.set(e,r),this},POP(e,t){let a=this.get(e);return a&&(a.count--,a.meta.delete(t)),this},SCROLL_PREVENT({doc:e,d:t,meta:a}){let r={doc:e,d:t,meta:La(a)},l=[$a(),Sa(),Da()];l.forEach(({before:n})=>n==null?void 0:n(r)),l.forEach(({after:n})=>n==null?void 0:n(r))},SCROLL_ALLOW({d:e}){e.dispose()},TEARDOWN({doc:e}){this.delete(e)}});Q.subscribe(()=>{let e=Q.getSnapshot(),t=new Map;for(let[a]of e)t.set(a,a.documentElement.style.overflow);for(let a of e.values()){let r=t.get(a.doc)==="hidden",l=a.count!==0;(l&&!r||!l&&r)&&Q.dispatch(a.count>0?"SCROLL_PREVENT":"SCROLL_ALLOW",a),a.count===0&&Q.dispatch("TEARDOWN",a)}});function Fa(e,t,a){let r=Ea(Q),l=m(()=>{let n=e.value?r.value.get(e.value):void 0;return n?n.count>0:!1});return Z([e,t],([n,o],[i],s)=>{if(!n||!o)return;Q.dispatch("PUSH",n,a);let u=!1;s(()=>{u||(Q.dispatch("POP",i??n,a),u=!0)})},{immediate:!0}),l}let Ae=new Map,se=new Map;function Qe(e,t=g(!0)){_(a=>{var r;if(!t.value)return;let l=S(e);if(!l)return;a(function(){var o;if(!l)return;let i=(o=se.get(l))!=null?o:1;if(i===1?se.delete(l):se.set(l,i-1),i!==1)return;let s=Ae.get(l);s&&(s["aria-hidden"]===null?l.removeAttribute("aria-hidden"):l.setAttribute("aria-hidden",s["aria-hidden"]),l.inert=s.inert,Ae.delete(l))});let n=(r=se.get(l))!=null?r:0;se.set(l,n+1),n===0&&(Ae.set(l,{"aria-hidden":l.getAttribute("aria-hidden"),inert:l.inert}),l.setAttribute("aria-hidden","true"),l.inert=!0)})}function Pa({defaultContainers:e=[],portals:t,mainTreeNodeRef:a}={}){let r=g(null),l=ee(r);function n(){var o,i,s;let u=[];for(let d of e)d!==null&&(d instanceof HTMLElement?u.push(d):"value"in d&&d.value instanceof HTMLElement&&u.push(d.value));if(t!=null&&t.value)for(let d of t.value)u.push(d);for(let d of(o=l==null?void 0:l.querySelectorAll("html > *, body > *"))!=null?o:[])d!==document.body&&d!==document.head&&d instanceof HTMLElement&&d.id!=="headlessui-portal-root"&&(d.contains(S(r))||d.contains((s=(i=S(r))==null?void 0:i.getRootNode())==null?void 0:s.host)||u.some(c=>d.contains(c))||u.push(d));return u}return{resolveContainers:n,contains(o){return n().some(i=>i.contains(o))},mainTreeNodeRef:r,MainTreeNode(){return a!=null?null:L(ke,{features:$e.Hidden,ref:r})}}}let wt=Symbol("ForcePortalRootContext");function Ca(){return R(wt,!1)}let Ze=A({name:"ForcePortalRoot",props:{as:{type:[Object,String],default:"template"},force:{type:Boolean,default:!1}},setup(e,{slots:t,attrs:a}){return I(wt,e.force),()=>{let{force:r,...l}=e;return x({theirProps:l,ourProps:{},slot:{},slots:t,attrs:a,name:"ForcePortalRoot"})}}}),Et=Symbol("StackContext");var Be=(e=>(e[e.Add=0]="Add",e[e.Remove=1]="Remove",e))(Be||{});function Oa(){return R(Et,()=>{})}function Aa({type:e,enabled:t,element:a,onUpdate:r}){let l=Oa();function n(...o){r==null||r(...o),l(...o)}O(()=>{Z(t,(o,i)=>{o?n(0,e,a):i===!0&&n(1,e,a)},{immediate:!0,flush:"sync"})}),k(()=>{t.value&&n(1,e,a)}),I(Et,n)}let Ra=Symbol("DescriptionContext");function _a({slot:e=g({}),name:t="Description",props:a={}}={}){let r=g([]);function l(n){return r.value.push(n),()=>{let o=r.value.indexOf(n);o!==-1&&r.value.splice(o,1)}}return I(Ra,{register:l,slot:e,name:t,props:a}),m(()=>r.value.length>0?r.value.join(" "):void 0)}function Na(e){let t=ee(e);if(!t){if(e===null)return null;throw new Error(`[Headless UI]: Cannot find ownerDocument for contextElement: ${e}`)}let a=t.getElementById("headlessui-portal-root");if(a)return a;let r=t.createElement("div");return r.setAttribute("id","headlessui-portal-root"),t.body.appendChild(r)}const He=new WeakMap;function Ma(e){var t;return(t=He.get(e))!=null?t:0}function et(e,t){let a=t(Ma(e));return a<=0?He.delete(e):He.set(e,a),a}let ka=A({name:"Portal",props:{as:{type:[Object,String],default:"div"}},setup(e,{slots:t,attrs:a}){let r=g(null),l=m(()=>ee(r)),n=Ca(),o=R(Tt,null),i=g(n===!0||o==null?Na(r.value):o.resolveTarget());i.value&&et(i.value,f=>f+1);let s=g(!1);O(()=>{s.value=!0}),_(()=>{n||o!=null&&(i.value=o.resolveTarget())});let u=R(Ie,null),d=!1,c=Ue();return Z(r,()=>{if(d||!u)return;let f=S(r);f&&(k(u.register(f),c),d=!0)}),k(()=>{var f,h;let v=(f=l.value)==null?void 0:f.getElementById("headlessui-portal-root");!v||i.value!==v||et(i.value,E=>E-1)||i.value.children.length>0||(h=i.value.parentElement)==null||h.removeChild(i.value)}),()=>{if(!s.value||i.value===null)return null;let f={ref:r,"data-headlessui-portal":""};return L(kt,{to:i.value},x({ourProps:f,theirProps:e,slot:{},attrs:a,slots:t,name:"Portal"}))}}}),Ie=Symbol("PortalParentContext");function Ba(){let e=R(Ie,null),t=g([]);function a(n){return t.value.push(n),e&&e.register(n),()=>r(n)}function r(n){let o=t.value.indexOf(n);o!==-1&&t.value.splice(o,1),e&&e.unregister(n)}let l={register:a,unregister:r,portals:t};return[t,A({name:"PortalWrapper",setup(n,{slots:o}){return I(Ie,l),()=>{var i;return(i=o.default)==null?void 0:i.call(o)}}})]}let Tt=Symbol("PortalGroupContext"),Ha=A({name:"PortalGroup",props:{as:{type:[Object,String],default:"template"},target:{type:Object,default:null}},setup(e,{attrs:t,slots:a}){let r=rt({resolveTarget(){return e.target}});return I(Tt,r),()=>{let{target:l,...n}=e;return x({theirProps:n,ourProps:{},slot:{},attrs:t,slots:a,name:"PortalGroup"})}}});var Ia=(e=>(e[e.Open=0]="Open",e[e.Closed=1]="Closed",e))(Ia||{});let xe=Symbol("DialogContext");function St(e){let t=R(xe,null);if(t===null){let a=new Error(`<${e} /> is missing a parent <Dialog /> component.`);throw Error.captureStackTrace&&Error.captureStackTrace(a,St),a}return t}let ge="DC8F892D-2EBD-447C-A4C8-A03058436FF4",xa=A({name:"Dialog",inheritAttrs:!1,props:{as:{type:[Object,String],default:"div"},static:{type:Boolean,default:!1},unmount:{type:Boolean,default:!0},open:{type:[Boolean,String],default:ge},initialFocus:{type:Object,default:null},id:{type:String,default:null},role:{type:String,default:"dialog"}},emits:{close:e=>!0},setup(e,{emit:t,attrs:a,slots:r,expose:l}){var n,o;let i=(n=e.id)!=null?n:`headlessui-dialog-${Ve()}`,s=g(!1);O(()=>{s.value=!0});let u=!1,d=m(()=>e.role==="dialog"||e.role==="alertdialog"?e.role:(u||(u=!0,console.warn(`Invalid role [${d}] passed to <Dialog />. Only \`dialog\` and and \`alertdialog\` are supported. Using \`dialog\` instead.`)),"dialog")),c=g(0),f=ot(),h=m(()=>e.open===ge&&f!==null?(f.value&C.Open)===C.Open:e.open),v=g(null),E=m(()=>ee(v));if(l({el:v,$el:v}),!(e.open!==ge||f!==null))throw new Error("You forgot to provide an `open` prop to the `Dialog`.");if(typeof h.value!="boolean")throw new Error(`You provided an \`open\` prop to the \`Dialog\`, but the value is not a boolean. Received: ${h.value===ge?void 0:e.open}`);let p=m(()=>s.value&&h.value?0:1),$=m(()=>p.value===0),b=m(()=>c.value>1),y=R(xe,null)!==null,[N,D]=Ba(),{resolveContainers:W,mainTreeNodeRef:z,MainTreeNode:j}=Pa({portals:N,defaultContainers:[m(()=>{var w;return(w=te.panelRef.value)!=null?w:v.value})]}),ve=m(()=>b.value?"parent":"leaf"),K=m(()=>f!==null?(f.value&C.Closing)===C.Closing:!1),Le=m(()=>y||K.value?!1:$.value),Fe=m(()=>{var w,T,F;return(F=Array.from((T=(w=E.value)==null?void 0:w.querySelectorAll("body > *"))!=null?T:[]).find(P=>P.id==="headlessui-portal-root"?!1:P.contains(S(z))&&P instanceof HTMLElement))!=null?F:null});Qe(Fe,Le);let M=m(()=>b.value?!0:$.value),oe=m(()=>{var w,T,F;return(F=Array.from((T=(w=E.value)==null?void 0:w.querySelectorAll("[data-headlessui-portal]"))!=null?T:[]).find(P=>P.contains(S(z))&&P instanceof HTMLElement))!=null?F:null});Qe(oe,M),Aa({type:"Dialog",enabled:m(()=>p.value===0),element:v,onUpdate:(w,T)=>{if(T==="Dialog")return q(w,{[Be.Add]:()=>c.value+=1,[Be.Remove]:()=>c.value-=1})}});let B=_a({name:"DialogDescription",slot:m(()=>({open:h.value}))}),H=g(null),te={titleId:H,panelRef:g(null),dialogState:p,setTitleId(w){H.value!==w&&(H.value=w)},close(){t("close",!1)}};I(xe,te);let Ke=m(()=>!(!$.value||b.value));pa(W,(w,T)=>{w.preventDefault(),te.close(),lt(()=>T==null?void 0:T.focus())},Ke);let Ge=m(()=>!(b.value||p.value!==0));ht((o=E.value)==null?void 0:o.defaultView,"keydown",w=>{Ge.value&&(w.defaultPrevented||w.key===Bt.Escape&&(w.preventDefault(),w.stopPropagation(),te.close()))});let Ye=m(()=>!(K.value||p.value!==0||y));return Fa(E,Ye,w=>{var T;return{containers:[...(T=w.containers)!=null?T:[],W]}}),_(w=>{if(p.value!==0)return;let T=S(v);if(!T)return;let F=new ResizeObserver(P=>{for(let Pe of P){let pe=Pe.target.getBoundingClientRect();pe.x===0&&pe.y===0&&pe.width===0&&pe.height===0&&te.close()}});F.observe(T),w(()=>F.disconnect())}),()=>{let{open:w,initialFocus:T,...F}=e,P={...a,ref:v,id:i,role:d.value,"aria-modal":p.value===0?!0:void 0,"aria-labelledby":H.value,"aria-describedby":B.value},Pe={open:p.value===0};return L(Ze,{force:!0},()=>[L(ka,()=>L(Ha,{target:v.value},()=>L(Ze,{force:!1},()=>L(ie,{initialFocus:T,containers:W,features:$.value?q(ve.value,{parent:ie.features.RestoreFocus,leaf:ie.features.All&~ie.features.FocusLock}):ie.features.None},()=>L(D,{},()=>x({ourProps:P,theirProps:{...F,...a},slot:Pe,attrs:a,slots:r,visible:p.value===0,features:_e.RenderStrategy|_e.Static,name:"Dialog"})))))),L(j)])}}}),ja=A({name:"DialogPanel",props:{as:{type:[Object,String],default:"div"},id:{type:String,default:null}},setup(e,{attrs:t,slots:a,expose:r}){var l;let n=(l=e.id)!=null?l:`headlessui-dialog-panel-${Ve()}`,o=St("DialogPanel");r({el:o.panelRef,$el:o.panelRef});function i(s){s.stopPropagation()}return()=>{let{...s}=e,u={id:n,ref:o.panelRef,onClick:i};return x({ourProps:u,theirProps:s,slot:{open:o.dialogState.value===0},attrs:t,slots:a,name:"DialogPanel"})}}});function Ua(e){let t={called:!1};return(...a)=>{if(!t.called)return t.called=!0,e(...a)}}function Re(e,...t){e&&t.length>0&&e.classList.add(...t)}function ye(e,...t){e&&t.length>0&&e.classList.remove(...t)}var je=(e=>(e.Finished="finished",e.Cancelled="cancelled",e))(je||{});function Va(e,t){let a=ce();if(!e)return a.dispose;let{transitionDuration:r,transitionDelay:l}=getComputedStyle(e),[n,o]=[r,l].map(i=>{let[s=0]=i.split(",").filter(Boolean).map(u=>u.includes("ms")?parseFloat(u):parseFloat(u)*1e3).sort((u,d)=>d-u);return s});return n!==0?a.setTimeout(()=>t("finished"),n+o):t("finished"),a.add(()=>t("cancelled")),a.dispose}function tt(e,t,a,r,l,n){let o=ce(),i=n!==void 0?Ua(n):()=>{};return ye(e,...l),Re(e,...t,...a),o.nextFrame(()=>{ye(e,...a),Re(e,...r),o.add(Va(e,s=>(ye(e,...r,...t),Re(e,...l),i(s))))}),o.add(()=>ye(e,...t,...a,...r,...l)),o.add(()=>i("cancelled")),o.dispose}function G(e=""){return e.split(/\s+/).filter(t=>t.length>1)}let We=Symbol("TransitionContext");var qa=(e=>(e.Visible="visible",e.Hidden="hidden",e))(qa||{});function Wa(){return R(We,null)!==null}function za(){let e=R(We,null);if(e===null)throw new Error("A <TransitionChild /> is used but it is missing a parent <TransitionRoot />.");return e}function Ka(){let e=R(ze,null);if(e===null)throw new Error("A <TransitionChild /> is used but it is missing a parent <TransitionRoot />.");return e}let ze=Symbol("NestingContext");function De(e){return"children"in e?De(e.children):e.value.filter(({state:t})=>t==="visible").length>0}function $t(e){let t=g([]),a=g(!1);O(()=>a.value=!0),k(()=>a.value=!1);function r(n,o=J.Hidden){let i=t.value.findIndex(({id:s})=>s===n);i!==-1&&(q(o,{[J.Unmount](){t.value.splice(i,1)},[J.Hidden](){t.value[i].state="hidden"}}),!De(t)&&a.value&&(e==null||e()))}function l(n){let o=t.value.find(({id:i})=>i===n);return o?o.state!=="visible"&&(o.state="visible"):t.value.push({id:n,state:"visible"}),()=>r(n,J.Unmount)}return{children:t,register:l,unregister:r}}let Dt=_e.RenderStrategy,Lt=A({props:{as:{type:[Object,String],default:"div"},show:{type:[Boolean],default:null},unmount:{type:[Boolean],default:!0},appear:{type:[Boolean],default:!1},enter:{type:[String],default:""},enterFrom:{type:[String],default:""},enterTo:{type:[String],default:""},entered:{type:[String],default:""},leave:{type:[String],default:""},leaveFrom:{type:[String],default:""},leaveTo:{type:[String],default:""}},emits:{beforeEnter:()=>!0,afterEnter:()=>!0,beforeLeave:()=>!0,afterLeave:()=>!0},setup(e,{emit:t,attrs:a,slots:r,expose:l}){let n=g(0);function o(){n.value|=C.Opening,t("beforeEnter")}function i(){n.value&=~C.Opening,t("afterEnter")}function s(){n.value|=C.Closing,t("beforeLeave")}function u(){n.value&=~C.Closing,t("afterLeave")}if(!Wa()&&Ht())return()=>L(Ft,{...e,onBeforeEnter:o,onAfterEnter:i,onBeforeLeave:s,onAfterLeave:u},r);let d=g(null),c=m(()=>e.unmount?J.Unmount:J.Hidden);l({el:d,$el:d});let{show:f,appear:h}=za(),{register:v,unregister:E}=Ka(),p=g(f.value?"visible":"hidden"),$={value:!0},b=Ve(),y={value:!1},N=$t(()=>{!y.value&&p.value!=="hidden"&&(p.value="hidden",E(b),u())});O(()=>{let M=v(b);k(M)}),_(()=>{if(c.value===J.Hidden&&b){if(f.value&&p.value!=="visible"){p.value="visible";return}q(p.value,{hidden:()=>E(b),visible:()=>v(b)})}});let D=G(e.enter),W=G(e.enterFrom),z=G(e.enterTo),j=G(e.entered),ve=G(e.leave),K=G(e.leaveFrom),Le=G(e.leaveTo);O(()=>{_(()=>{if(p.value==="visible"){let M=S(d);if(M instanceof Comment&&M.data==="")throw new Error("Did you forget to passthrough the `ref` to the actual DOM node?")}})});function Fe(M){let oe=$.value&&!h.value,B=S(d);!B||!(B instanceof HTMLElement)||oe||(y.value=!0,f.value&&o(),f.value||s(),M(f.value?tt(B,D,W,z,j,H=>{y.value=!1,H===je.Finished&&i()}):tt(B,ve,K,Le,j,H=>{y.value=!1,H===je.Finished&&(De(N)||(p.value="hidden",E(b),u()))})))}return O(()=>{Z([f],(M,oe,B)=>{Fe(B),$.value=!1},{immediate:!0})}),I(ze,N),It(m(()=>q(p.value,{visible:C.Open,hidden:C.Closed})|n.value)),()=>{let{appear:M,show:oe,enter:B,enterFrom:H,enterTo:te,entered:Ke,leave:Ge,leaveFrom:Ye,leaveTo:w,...T}=e,F={ref:d},P={...T,...h.value&&f.value&&fe.isServer?{class:le([a.class,T.class,...D,...W])}:{}};return x({theirProps:P,ourProps:F,slot:{},slots:r,attrs:a,features:Dt,visible:p.value==="visible",name:"TransitionChild"})}}}),Ga=Lt,Ft=A({inheritAttrs:!1,props:{as:{type:[Object,String],default:"div"},show:{type:[Boolean],default:null},unmount:{type:[Boolean],default:!0},appear:{type:[Boolean],default:!1},enter:{type:[String],default:""},enterFrom:{type:[String],default:""},enterTo:{type:[String],default:""},entered:{type:[String],default:""},leave:{type:[String],default:""},leaveFrom:{type:[String],default:""},leaveTo:{type:[String],default:""}},emits:{beforeEnter:()=>!0,afterEnter:()=>!0,beforeLeave:()=>!0,afterLeave:()=>!0},setup(e,{emit:t,attrs:a,slots:r}){let l=ot(),n=m(()=>e.show===null&&l!==null?(l.value&C.Open)===C.Open:e.show);_(()=>{if(![!0,!1].includes(n.value))throw new Error('A <Transition /> is used but it is missing a `:show="true | false"` prop.')});let o=g(n.value?"visible":"hidden"),i=$t(()=>{o.value="hidden"}),s=g(!0),u={show:n,appear:m(()=>e.appear||!s.value)};return O(()=>{_(()=>{s.value=!1,n.value?o.value="visible":De(i)||(o.value="hidden")})}),I(ze,i),I(We,u),()=>{let d=xt(e,["show","appear","unmount","onBeforeEnter","onBeforeLeave","onAfterEnter","onAfterLeave"]),c={unmount:e.unmount};return x({ourProps:{...c,as:"template"},theirProps:{},slot:{},slots:{...r,default:()=>[L(Ga,{onBeforeEnter:()=>t("beforeEnter"),onAfterEnter:()=>t("afterEnter"),onBeforeLeave:()=>t("beforeLeave"),onAfterLeave:()=>t("afterLeave"),...a,...c,...d},r.default)]},attrs:{},features:Dt,visible:o.value==="visible",name:"Transition"})}}});const Ya=st(Se.ui.strategy,Se.ui.card,aa),Xa=A({inheritAttrs:!1,props:{as:{type:String,default:"div"},class:{type:[String,Object,Array],default:()=>""},ui:{type:Object,default:()=>({})}},setup(e){const{ui:t,attrs:a}=ut("card",Te(e,"ui"),Ya),r=m(()=>jt(Ut(t.value.base,t.value.rounded,t.value.divide,t.value.ring,t.value.shadow,t.value.background),e.class));return{ui:t,attrs:a,cardClass:r}}});function Ja(e,t,a,r,l,n){return re(),Ne(Vt(e.$attrs.onSubmit?"form":e.as),Ee({class:e.cardClass},e.attrs),{default:ae(()=>[e.$slots.header?(re(),Ce("div",{key:0,class:le([e.ui.header.base,e.ui.header.padding,e.ui.header.background])},[be(e.$slots,"header")],2)):we("",!0),e.$slots.default?(re(),Ce("div",{key:1,class:le([e.ui.body.base,e.ui.body.padding,e.ui.body.background])},[be(e.$slots,"default")],2)):we("",!0),e.$slots.footer?(re(),Ce("div",{key:2,class:le([e.ui.footer.base,e.ui.footer.padding,e.ui.footer.background])},[be(e.$slots,"footer")],2)):we("",!0)]),_:3},16,["class"])}const ln=it(Xa,[["render",Ja]]),Qa=st(Se.ui.strategy,Se.ui.slideover,na),Za=A({components:{HDialog:xa,HDialogPanel:ja,TransitionRoot:Ft,TransitionChild:Lt},inheritAttrs:!1,props:{modelValue:{type:Boolean,default:!1},appear:{type:Boolean,default:!1},side:{type:String,default:"right",validator:e=>["left","right","top","bottom"].includes(e)},overlay:{type:Boolean,default:!0},transition:{type:Boolean,default:!0},preventClose:{type:Boolean,default:!1},class:{type:[String,Object,Array],default:()=>""},ui:{type:Object,default:()=>({})}},emits:["update:modelValue","close","close-prevented","after-leave"],setup(e,{emit:t}){const{ui:a,attrs:r}=ut("slideover",Te(e,"ui"),Qa,Te(e,"class")),l=m({get(){return e.modelValue},set(u){t("update:modelValue",u)}}),n=m(()=>{if(!e.transition)return{};let u,d;switch(e.side){case"left":u=a.value.translate.left,d=a.value.translate.left;break;case"right":u=a.value.translate.right,d=a.value.translate.right;break;case"top":u=a.value.translate.top,d=a.value.translate.top;break;case"bottom":u=a.value.translate.bottom,d=a.value.translate.bottom;break;default:u=a.value.translate.right,d=a.value.translate.right}return{...a.value.transition,enterFrom:u,enterTo:a.value.translate.base,leaveFrom:a.value.translate.base,leaveTo:d}}),o=m(()=>{switch(e.side){case"left":return"horizontal";case"right":return"horizontal";case"top":return"vertical";case"bottom":return"vertical";default:return"right"}});function i(u){if(e.preventClose){t("close-prevented");return}l.value=u,t("close")}const s=()=>{t("after-leave")};return qt(()=>ta("$zSEivjycv8")),{ui:a,attrs:r,isOpen:l,transitionClass:n,sideType:o,onAfterLeave:s,close:i}}});function en(e,t,a,r,l,n){const o=me("TransitionChild"),i=me("HDialogPanel"),s=me("HDialog"),u=me("TransitionRoot");return re(),Ne(u,{as:"template",appear:e.appear,show:e.isOpen,onAfterLeave:e.onAfterLeave},{default:ae(()=>[Oe(s,Ee({class:[e.ui.wrapper,{"justify-end":e.side==="right"},{"items-end":e.side==="bottom"}]},e.attrs,{onClose:e.close}),{default:ae(()=>[e.overlay?(re(),Ne(o,Ee({key:0,as:"template",appear:e.appear},e.ui.overlay.transition),{default:ae(()=>[Wt("div",{class:le([e.ui.overlay.base,e.ui.overlay.background])},null,2)]),_:1},16,["appear"])):we("",!0),Oe(o,Ee({as:"template",appear:e.appear},e.transitionClass),{default:ae(()=>[Oe(i,{class:le([e.ui.base,e.sideType==="horizontal"?[e.ui.width,"h-full"]:[e.ui.height,"w-full"],e.ui.background,e.ui.ring,e.ui.padding])},{default:ae(()=>[be(e.$slots,"default")]),_:3},8,["class"])]),_:3},16,["appear"])]),_:3},16,["class","onClose"])]),_:3},8,["appear","show","onAfterLeave"])}const on=it(Za,[["render",en]]),sn=(e,t={})=>{const a=zt("/go-api",e.toString());return Jt(a,{server:!1,...t},"$sfBywraVUZ")};export{U as N,ca as O,ln as _,vt as a,nn as b,$e as c,ta as d,on as e,ke as f,ft as h,ee as i,va as n,ce as o,X as t,sn as u,rn as v,pa as w};
