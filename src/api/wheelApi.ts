import httpClient from './httpClient';
import type { WheelEntry, WheelSettings } from '../types/wheel';

export interface CreateWheelPayload {
  title: string;
  entries: Omit<WheelEntry, 'id' | 'position'>[];
  settings: WheelSettings;
}

export interface UpdateWheelPayload {
  title: string;
  entries: { label: string; weight: number; color?: string; position: number }[];
  settings: WheelSettings;
}

export const wheelApi = {
  health() {
    return httpClient.get('/health');
  },

  createWheel(payload: CreateWheelPayload) {
    return httpClient.post('/wheels', payload);
  },

  getWheel(shareCode: string) {
    return httpClient.get(`/wheels/${shareCode}`);
  },

  updateWheel(id: string, payload: UpdateWheelPayload) {
    return httpClient.put(`/wheels/${id}`, payload);
  },

  deleteWheel(id: string) {
    return httpClient.delete(`/wheels/${id}`);
  },

  recordSpin(id: string, entryId?: string, resultLabel?: string) {
    return httpClient.post(`/wheels/${id}/spin`, { entryId, resultLabel });
  },

  getHistory(id: string) {
    return httpClient.get(`/wheels/${id}/history`);
  },

  clearHistory(id: string) {
    return httpClient.delete(`/wheels/${id}/history`);
  },

  duplicateWheel(id: string) {
    return httpClient.post(`/wheels/${id}/duplicate`);
  }
};
