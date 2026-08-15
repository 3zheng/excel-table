// src/constants/index.ts

/** 地区列表（value 用中文，与后端、数据库保持一致） */
export const REGIONS = [
  { value: '玻利维亚', labelKey: 'bolivia' },
  { value: '秘鲁', labelKey: 'peru' },
  { value: '智利', labelKey: 'chile' },
  { value: '西班牙', labelKey: 'spain' },
  { value: '乌拉圭', labelKey: 'uruguay' },
  { value: '委内瑞拉', labelKey: 'venezuela' },
  { value: '阿根廷', labelKey: 'argentina' },
  { value: '巴拉圭', labelKey: 'paraguay' },
  { value: '巴拿马', labelKey: 'panama' },
  { value: '日本', labelKey: 'japan' },
] as const

/** 买方列表（临时数据，后续替换为真实买方） */
export interface Buyer {
  name: string
  address: string
  tel: string
}

export const BUYERS: Buyer[] = [
  {
    name: '买方A',
    address: '地址A示例',
    tel: '电话A示例',
  },
  {
    name: '买方B',
    address: '地址B示例',
    tel: '电话B示例',
  },
  {
    name: '买方C',
    address: '地址C示例',
    tel: '电话C示例',
  },
]