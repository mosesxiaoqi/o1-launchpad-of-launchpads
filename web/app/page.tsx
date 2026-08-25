import Link from 'next/link'

import { SiteHeader } from '@/components/site-header'

export default function Home() {
  return (
    <main className="site-shell home-page">
      <SiteHeader />
      <section className="hero">
        <div className="hero-copy">
          <p className="eyebrow"><span>01</span> MULTI-TENANT TOKEN INFRASTRUCTURE</p>
          <h1>发行你的 Launchpad。<br /><em>掌控每一次启动。</em></h1>
          <p className="hero-lede">在一套经过验证的链上基础设施上建立独立品牌、发行固定供应量 Token，并让费用自动流向正确的参与者。</p>
          <div className="hero-actions">
            <Link className="button button-primary" href="/create-launchpad">创建 Launchpad <span aria-hidden="true">↗</span></Link>
            <a className="text-link" href="#protocol">查看协议参数 <span aria-hidden="true">↓</span></a>
          </div>
        </div>
        <aside className="protocol-console" aria-label="协议状态">
          <div className="console-header"><span>PROTOCOL / STATUS</span><span className="live-indicator"><i /> LIVE</span></div>
          <div className="console-core"><div className="orbit"><span>O1</span></div><p>NETWORK</p><strong>Base Sepolia</strong></div>
          <dl className="console-list">
            <div><dt>CHAIN ID</dt><dd>84532</dd></div>
            <div><dt>QUOTE ASSET</dt><dd>NATIVE ETH</dd></div>
            <div><dt>POOL ENGINE</dt><dd>UNISWAP V4</dd></div>
            <div><dt>STATUS</dt><dd className="status-ok">OPERATIONAL</dd></div>
          </dl>
        </aside>
      </section>
      <section className="metric-strip" id="protocol" aria-label="协议参数">
        <Metric index="01" value="1,000,000,000" label="固定供应量" />
        <Metric index="02" value="1.5%" label="正常阶段总费率" />
        <Metric index="03" value="16 秒" label="ANTI-SNIPE 保护" />
        <Metric index="04" value="∞" label="永久流动性" />
      </section>
      <section className="workflow-section">
        <div className="section-heading">
          <p className="eyebrow"><span>02</span> HOW IT WORKS</p>
          <h2>从品牌到市场，<br />一条清晰的链上路径。</h2>
        </div>
        <div className="workflow-grid">
          <Workflow number="01" title="建立品牌" copy="注册独立 Launchpad、链上 Treasury 与品牌标识。" tag="REGISTRY" />
          <Workflow number="02" title="发行资产" copy="创建固定供应量 ERC-20，并自动初始化永久流动性。" tag="FACTORY" />
          <Workflow number="03" title="开放交易" copy="通过 Uniswap v4 Hook 执行 Swap、费用拆分与保护机制。" tag="MARKET" />
        </div>
      </section>
    </main>
  )
}

function Metric({ index, value, label }: { index: string; value: string; label: string }) {
  return <article><small>{index}</small><strong>{value}</strong><span>{label}</span></article>
}

function Workflow({ number, title, copy, tag }: { number: string; title: string; copy: string; tag: string }) {
  return <article className="workflow-card"><div><span>{number}</span><small>{tag}</small></div><h3>{title}</h3><p>{copy}</p></article>
}
