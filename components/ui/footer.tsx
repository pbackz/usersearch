import * as React from 'react';

const Footer: React.FC = () => {
  return (
    <footer className="py-4 text-center text-sm text-muted-foreground">
      &copy; {new Date().getFullYear()} Pierre Baconnier. All rights reserved.
    </footer>
  );
};

export default Footer;