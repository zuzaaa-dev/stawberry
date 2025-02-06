import React from 'react';
import { Link } from 'react-router-dom';

const Dashboard = () => {
  return (
    <div>
      <h1>Dashboard</h1>
      <nav>
        <ul>
          <li><Link to="/products">Products</Link></li>
          <li><Link to="/offers">Offers</Link></li>
          <li><Link to="/health">Health Status</Link></li>
        </ul>
      </nav>
    </div>
  );
};

export default Dashboard;