import { useState, useEffect } from 'react';
import { featureToggleService } from '../services/featureToggle';

export const useFeatureFlag = (flagName: string) => {
  const [isEnabled, setIsEnabled] = useState<boolean>(false);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    const checkFeatureFlag = async () => {
      try {
        setIsLoading(true);
        const enabled = await featureToggleService.isEnabled(flagName);
        setIsEnabled(enabled);
        setError(null);
      } catch (err) {
        setError(err instanceof Error ? err : new Error('Failed to check feature flag'));
      } finally {
        setIsLoading(false);
      }
    };

    checkFeatureFlag();
  }, [flagName]);

  const toggleFlag = async (enabled: boolean) => {
    try {
      await featureToggleService.setFlag(flagName, enabled);
      setIsEnabled(enabled);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err : new Error('Failed to toggle feature flag'));
      throw err;
    }
  };

  return { isEnabled, isLoading, error, toggleFlag };
}; 