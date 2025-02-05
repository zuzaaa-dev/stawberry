import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { Link } from 'react-router-dom';

const OfferList = () => {
  const [offers, setOffers] = useState([]);
  const [error, setError] = useState(null);

  useEffect(() => {
    axios.get('/api/offers')
      .then(response => {
        setOffers(response.data.data);
      })
      .catch(error => {
        setError(error.message);
      });
  }, []);

  if (error) {
    return <div>Error: {error}</div>;
  }

  return (
    <div>
      <h1>Offers</h1>
      <ul>
        {offers.map(offer => (
          <li key={offer.id}>
            <Link to={`/offers/${offer.id}`}>Offer for Product {offer.product_id}</Link>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default OfferList;