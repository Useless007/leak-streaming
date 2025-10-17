import { render, screen } from '@testing-library/react';
import HomePage from '../page';

describe('HomePage', () => {
  it('แสดงข้อความหลักของหน้า landing', () => {
    render(<HomePage />);
    expect(screen.getByText(/Leak Streaming Studio/i)).toBeInTheDocument();
    expect(
      screen.getByRole('link', {
        name: /เริ่มสำรวจหนัง/i
      })
    ).toBeInTheDocument();
  });
});
