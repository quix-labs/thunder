import{e as d,A as h,C as g,bd as b,h as w,o as i,i as v,w as r,b as y,j as l,be as z,aL as C,aM as B,y as p,c as k,Q as x,d as V,t as M,m as j,bf as A,z as c,af as N}from"./DJAyu_pX.js";const _={slots:{root:"relative inline-flex items-center justify-center shrink-0",base:"rounded-full ring ring-[var(--ui-bg)] flex items-center justify-center text-[var(--ui-bg)] font-medium whitespace-nowrap"},variants:{color:{primary:"bg-[var(--ui-primary)]",secondary:"bg-[var(--ui-secondary)]",success:"bg-[var(--ui-success)]",info:"bg-[var(--ui-info)]",warning:"bg-[var(--ui-warning)]",error:"bg-[var(--ui-error)]",neutral:"bg-[var(--ui-text-muted)]"},size:{"3xs":"h-[4px] min-w-[4px] text-[4px]","2xs":"h-[5px] min-w-[5px] text-[5px]",xs:"h-[6px] min-w-[6px] text-[6px]",sm:"h-[7px] min-w-[7px] text-[7px]",md:"h-[8px] min-w-[8px] text-[8px]",lg:"h-[9px] min-w-[9px] text-[9px]",xl:"h-[10px] min-w-[10px] text-[10px]","2xl":"h-[11px] min-w-[11px] text-[11px]","3xl":"h-[12px] min-w-[12px] text-[12px]"},position:{"top-right":"top-0 right-0","bottom-right":"bottom-0 right-0","top-left":"top-0 left-0","bottom-left":"bottom-0 left-0"},inset:{false:""},standalone:{false:"absolute"}},compoundVariants:[{position:"top-right",inset:!1,class:"-translate-y-1/2 translate-x-1/2 transform"},{position:"bottom-right",inset:!1,class:"translate-y-1/2 translate-x-1/2 transform"},{position:"top-left",inset:!1,class:"-translate-y-1/2 -translate-x-1/2 transform"},{position:"bottom-left",inset:!1,class:"translate-y-1/2 -translate-x-1/2 transform"}],defaultVariants:{size:"md",color:"primary",position:"top-right"}},$=N;var u;const P=c({extend:c(_),...((u=$.ui)==null?void 0:u.chip)||{}}),E=d({inheritAttrs:!1,__name:"Chip",props:h({as:{},text:{},color:{},size:{},position:{},inset:{type:Boolean,default:!1},standalone:{type:Boolean,default:!1},class:{},ui:{}},{show:{type:Boolean,default:!0},showModifiers:{}}),emits:["update:show"],setup(s){const t=s,m=g(s,"show"),{size:f}=b(t),a=w(()=>P({color:t.color,size:f.value,position:t.position,inset:t.inset,standalone:t.standalone}));return(e,S)=>{var o;return i(),v(l(A),{as:e.as,class:x(a.value.root({class:[t.class,(o=t.ui)==null?void 0:o.root]}))},{default:r(()=>{var n;return[y(l(z),C(B(e.$attrs)),{default:r(()=>[p(e.$slots,"default")]),_:3},16),m.value?(i(),k("span",{key:0,class:x(a.value.base({class:(n=t.ui)==null?void 0:n.base}))},[p(e.$slots,"content",{},()=>[V(M(e.text),1)])],2)):j("",!0)]}),_:3},8,["as","class"])}}});export{E as _};
