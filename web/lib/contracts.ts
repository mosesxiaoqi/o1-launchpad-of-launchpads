import { parseAbi } from 'viem'

export const launchpadRegistryAbi = parseAbi([
  'function createLaunchpad(string normalizedSlug, address treasury) returns (bytes32 id)',
  'event LaunchpadCreated(bytes32 indexed id, address indexed owner, address indexed treasury, string slug)',
])

export const multiTenantFactoryAbi = parseAbi([
  'function launch((bytes32 launchpadId,string name,string symbol,string contractURI,bytes32 salt,address quote,uint64 expectedConfigVersion,uint64 deadline) params) returns (address token, bytes32 poolId)',
  'event Launched(address indexed token, bytes32 indexed poolId, bytes32 indexed launchpadId, address creator, address quote, uint256 supply, int24 tickSpacing)',
])

export const feeEscrowAbi = parseAbi([
  'function owed(address recipient, address currency) view returns (uint256)',
  'function claim(address recipient, address currency)',
])
