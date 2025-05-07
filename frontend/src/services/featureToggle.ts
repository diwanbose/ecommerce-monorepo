import axios from 'axios';
import { FeatureFlag } from '../types';

const API_URL = process.env.REACT_APP_FEATURE_TOGGLE_URL || 'http://localhost:8084';

export const featureToggleService = {
  async isEnabled(flagName: string): Promise<boolean> {
    try {
      const response = await axios.get<{ enabled: boolean }>(`${API_URL}/api/flags/${flagName}`);
      return response.data.enabled;
    } catch (error) {
      console.error(`Error checking feature flag ${flagName}:`, error);
      return false;
    }
  },

  async setFlag(flagName: string, enabled: boolean): Promise<void> {
    try {
      await axios.post(`${API_URL}/api/flags/${flagName}`, { enabled });
    } catch (error) {
      console.error(`Error setting feature flag ${flagName}:`, error);
      throw error;
    }
  },

  async getAllFlags(): Promise<FeatureFlag[]> {
    try {
      const response = await axios.get<FeatureFlag[]>(`${API_URL}/api/flags`);
      return response.data;
    } catch (error) {
      console.error('Error fetching feature flags:', error);
      return [];
    }
  }
}; 