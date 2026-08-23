'use client'

import { useRouter } from 'next/navigation'
import { useState, type FormEvent, type ReactNode } from 'react'
import { useAccount, usePublicClient, useSignMessage, useWriteContract } from 'wagmi'

import { launchpadApi } from '@/lib/api'
import { launchpadRegistryAbi } from '@/lib/contracts'
import { normalizeSlug } from '@/lib/launchpad-id'
import { chain } from '@/lib/wagmi'
import { TransactionStatus, type TransactionPhase } from './transaction-status'

type Fields = {
  name: string
  slug: string
  description: string
  logoUrl: string
  primaryColor: string
}

type FieldErrors = Partial<Record<keyof Fields, string>>

const initialFields: Fields = {
  name: '',
  slug: '',
  description: '',
  logoUrl: '',
  primaryColor: '#5B5CF6',
}

function validate(fields: Fields): { slug?: string; errors: FieldErrors } {
  const errors: FieldErrors = {}
  let slug: string | undefined
  try {
    slug = normalizeSlug(fields.slug)
  } catch {
    errors.slug = 'Slug 需要为 3–32 位小写字母、数字或内部单个连字符'
  }
  if (!fields.name.trim() || fields.name.trim().length > 80) errors.name = '名称需要为 1–80 个字符'
  if (new TextEncoder().encode(fields.description).length > 1000) errors.description = '简介不能超过 1000 bytes'
  if (fields.logoUrl) {
    try {
      if (new URL(fields.logoUrl).protocol !== 'https:') throw new Error()
    } catch {
      errors.logoUrl = 'Logo 地址必须是 HTTPS URL'
    }
  }
  if (!/^#[0-9a-fA-F]{6}$/.test(fields.primaryColor)) errors.primaryColor = '主题色格式应为 #RRGGBB'
  return { slug, errors }
}

export function LaunchpadForm() {
  const router = useRouter()
  const { address, chainId, isConnected } = useAccount()
  const { signMessageAsync } = useSignMessage()
  const { writeContractAsync } = useWriteContract()
  const publicClient = usePublicClient()
  const [fields, setFields] = useState(initialFields)
  const [fieldErrors, setFieldErrors] = useState<FieldErrors>({})
  const [phase, setPhase] = useState<TransactionPhase>('idle')
  const [error, setError] = useState('')
  const busy = phase !== 'idle' && phase !== 'success'
  const firstError = (['slug', 'name', 'description', 'logoUrl', 'primaryColor'] as const).find(
    (field) => fieldErrors[field],
  )

  function update(field: keyof Fields, value: string) {
    setFields((current) => ({ ...current, [field]: value }))
    setFieldErrors((current) => ({ ...current, [field]: undefined }))
  }

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    if (!address || !isConnected || chainId !== chain.id || !publicClient) return
    const checked = validate(fields)
    if (!checked.slug || Object.keys(checked.errors).length) {
      setFieldErrors(checked.errors)
      return
    }

    setError('')
    try {
      setPhase('auth')
      const challenge = await launchpadApi.createChallenge({
        chain_id: chain.id,
        address,
        domain: window.location.hostname,
        uri: window.location.origin,
      })
      const signature = await signMessageAsync({ message: challenge.message })
      await launchpadApi.verifySignature({ challenge_id: challenge.challenge_id, address, signature })

      const config = await launchpadApi.config()
      setPhase('wallet')
      const hash = await writeContractAsync({
        address: config.registry,
        abi: launchpadRegistryAbi,
        functionName: 'createLaunchpad',
        args: [checked.slug, address],
        chainId: chain.id,
      })

      setPhase('confirming')
      const receipt = await publicClient.waitForTransactionReceipt({ hash, confirmations: 2 })
      if (receipt.status !== 'success') throw new Error('Registry 交易已回滚')

      setPhase('saving')
      await launchpadApi.createLaunchpad({
        chain_id: chain.id,
        slug: checked.slug,
        name: fields.name.trim(),
        description: fields.description,
        logo_url: fields.logoUrl,
        primary_color: fields.primaryColor,
        registry_tx_hash: hash,
      })
      setPhase('success')
      router.push(`/launchpad/${checked.slug}`)
    } catch (cause) {
      setPhase('idle')
      setError(cause instanceof Error ? cause.message : '创建 Launchpad 失败')
    }
  }

  const availability = !isConnected
    ? '请先连接钱包'
    : chainId !== chain.id
      ? '请切换到 Base Sepolia'
      : ''

  return (
    <form onSubmit={submit} noValidate>
      <p>{availability}</p>
      <Field label="Slug" error={fieldErrors.slug}>
        <input ref={firstError === 'slug' ? focusElement : undefined} id="slug" value={fields.slug} onChange={(event) => update('slug', event.target.value)} aria-describedby={fieldErrors.slug ? 'slug-error' : undefined} />
      </Field>
      <Field label="名称" error={fieldErrors.name}>
        <input ref={firstError === 'name' ? focusElement : undefined} id="name" value={fields.name} onChange={(event) => update('name', event.target.value)} aria-describedby={fieldErrors.name ? 'name-error' : undefined} />
      </Field>
      <Field label="简介" error={fieldErrors.description}>
        <textarea ref={firstError === 'description' ? focusElement : undefined} id="description" value={fields.description} onChange={(event) => update('description', event.target.value)} aria-describedby={fieldErrors.description ? 'description-error' : undefined} />
      </Field>
      <Field label="Logo URL" error={fieldErrors.logoUrl}>
        <input ref={firstError === 'logoUrl' ? focusElement : undefined} id="logoUrl" type="url" value={fields.logoUrl} onChange={(event) => update('logoUrl', event.target.value)} aria-describedby={fieldErrors.logoUrl ? 'logoUrl-error' : undefined} />
      </Field>
      <Field label="主题色" error={fieldErrors.primaryColor}>
        <input ref={firstError === 'primaryColor' ? focusElement : undefined} id="primaryColor" value={fields.primaryColor} onChange={(event) => update('primaryColor', event.target.value)} aria-describedby={fieldErrors.primaryColor ? 'primaryColor-error' : undefined} />
      </Field>
      <button type="submit" disabled={Boolean(availability) || busy}>创建 Launchpad</button>
      <TransactionStatus phase={phase} error={error} />
    </form>
  )
}

function Field({ label, error, children }: { label: keyof typeof fieldIDs; error?: string; children: ReactNode }) {
  const id = fieldIDs[label]
  return (
    <div>
      <label htmlFor={id}>{label}</label>
      {children}
      {error && <p id={`${id}-error`}>{error}</p>}
    </div>
  )
}

const fieldIDs = { Slug: 'slug', 名称: 'name', 简介: 'description', 'Logo URL': 'logoUrl', 主题色: 'primaryColor' } as const

function focusElement(element: HTMLElement | null) {
  element?.focus()
}
