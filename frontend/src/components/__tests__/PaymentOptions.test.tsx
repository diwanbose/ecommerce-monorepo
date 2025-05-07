import React from 'react';
import { screen, fireEvent } from '@testing-library/react';
import { render } from '../../test-utils';
import { PaymentOptions } from '../PaymentOptions';
import { useFeatureFlag } from '../../hooks/useFeatureFlag';

// Mock the useFeatureFlag hook
jest.mock('../../hooks/useFeatureFlag');

describe('PaymentOptions', () => {
  const mockOnPaymentChange = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders all payment options when COD is enabled', () => {
    (useFeatureFlag as jest.Mock).mockReturnValue({
      isEnabled: true,
      isLoading: false,
      error: null,
    });

    render(
      <PaymentOptions
        selectedPayment="credit_card"
        onPaymentChange={mockOnPaymentChange}
      />
    );

    expect(screen.getByLabelText('Credit Card')).toBeInTheDocument();
    expect(screen.getByLabelText('Net Banking')).toBeInTheDocument();
    expect(screen.getByLabelText('Cash on Delivery')).toBeInTheDocument();
  });

  it('does not render COD option when feature flag is disabled', () => {
    (useFeatureFlag as jest.Mock).mockReturnValue({
      isEnabled: false,
      isLoading: false,
      error: null,
    });

    render(
      <PaymentOptions
        selectedPayment="credit_card"
        onPaymentChange={mockOnPaymentChange}
      />
    );

    expect(screen.getByLabelText('Credit Card')).toBeInTheDocument();
    expect(screen.getByLabelText('Net Banking')).toBeInTheDocument();
    expect(screen.queryByLabelText('Cash on Delivery')).not.toBeInTheDocument();
  });

  it('shows loading state when feature flag is loading', () => {
    (useFeatureFlag as jest.Mock).mockReturnValue({
      isEnabled: false,
      isLoading: true,
      error: null,
    });

    render(
      <PaymentOptions
        selectedPayment="credit_card"
        onPaymentChange={mockOnPaymentChange}
      />
    );

    expect(screen.getByText('Loading payment options...')).toBeInTheDocument();
  });

  it('calls onPaymentChange when a payment option is selected', () => {
    (useFeatureFlag as jest.Mock).mockReturnValue({
      isEnabled: true,
      isLoading: false,
      error: null,
    });

    render(
      <PaymentOptions
        selectedPayment="credit_card"
        onPaymentChange={mockOnPaymentChange}
      />
    );

    fireEvent.click(screen.getByLabelText('Net Banking'));
    expect(mockOnPaymentChange).toHaveBeenCalledWith('netbanking');
  });
}); 