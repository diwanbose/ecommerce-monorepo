import React from 'react';
import { Box, FormControl, FormLabel, RadioGroup, FormControlLabel, Radio } from '@mui/material';
import { useFeatureFlag } from '../hooks/useFeatureFlag';

interface PaymentOptionsProps {
  selectedPayment: string;
  onPaymentChange: (payment: string) => void;
}

export const PaymentOptions: React.FC<PaymentOptionsProps> = ({
  selectedPayment,
  onPaymentChange,
}) => {
  const { isEnabled: isCodEnabled, isLoading } = useFeatureFlag('enableCodPayment');

  if (isLoading) {
    return <Box>Loading payment options...</Box>;
  }

  return (
    <FormControl component="fieldset">
      <FormLabel component="legend">Payment Method</FormLabel>
      <RadioGroup
        value={selectedPayment}
        onChange={(e) => onPaymentChange(e.target.value)}
      >
        <FormControlLabel
          value="credit_card"
          control={<Radio />}
          label="Credit Card"
        />
        <FormControlLabel
          value="netbanking"
          control={<Radio />}
          label="Net Banking"
        />
        {isCodEnabled && (
          <FormControlLabel
            value="cod"
            control={<Radio />}
            label="Cash on Delivery"
          />
        )}
      </RadioGroup>
    </FormControl>
  );
}; 