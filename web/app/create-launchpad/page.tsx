import { LaunchpadForm } from '@/components/launchpad-form'
import { SiteHeader } from '@/components/site-header'

export default function CreateLaunchpadPage() {
  return (
    <main className="site-shell app-page">
      <SiteHeader />
      <section className="page-intro">
        <p className="eyebrow"><span>CREATE / 01</span> TENANT REGISTRY</p>
        <h1>建立你的<br /><em>发行品牌。</em></h1>
        <p>一次链上注册，获得独立 Treasury、品牌空间与 Token 发行入口。</p>
      </section>
      <div className="form-layout">
        <aside className="context-panel">
          <p className="panel-label">BEFORE YOU START</p>
          <ol className="numbered-list">
            <li><span>01</span><div><strong>连接钱包</strong><p>钱包将成为 Launchpad 的链上 Owner。</p></div></li>
            <li><span>02</span><div><strong>注册身份</strong><p>Slug 与 Treasury 写入 Base Sepolia Registry。</p></div></li>
            <li><span>03</span><div><strong>保存品牌</strong><p>名称、简介和视觉配置保存在平台数据库。</p></div></li>
          </ol>
          <div className="notice"><span>!</span><p>测试网 MVP。请勿使用生产资金或主网钱包。</p></div>
        </aside>
        <div className="form-panel"><LaunchpadForm /></div>
      </div>
    </main>
  )
}
