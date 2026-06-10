export interface WheelEntry {
  id: string;
  label: string;
  weight: number;
  color?: string;
  position: number;
}

export interface WheelSettings {
  theme: 'default' | 'neon' | 'candy' | 'pastel' | 'monochrome' | 'ocean' | 'custom';
  customColors: string[];
  spinDuration: number; // In seconds (e.g., 3, 5, 8)
  volume: number; // 0 - 100
  enableSound: boolean;
  enableTickSound: boolean;
  enableVictorySound: boolean;
  enableConfetti: boolean;
  allowDuplicates: boolean;
  autoRemoveWinner: boolean;
  enableWeights: boolean;
  backgroundImage?: string;
  logoUrl?: string;
  showTextOnWheel: boolean;
  fontSize: number; // E.g., 14, 18, 24
}

export interface SpinHistory {
  id: string;
  resultLabel: string;
  spunAt: string; // ISO String
  entryId?: string;
}

export interface Wheel {
  id: string;
  title: string;
  shareCode: string;
  settings: WheelSettings;
  permission: 'view' | 'spin' | 'edit';
  entries: WheelEntry[];
  createdAt?: string;
  updatedAt?: string;
}
