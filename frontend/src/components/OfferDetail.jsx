import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useParams } from 'react-router-dom';

const OfferDetail = () => {
  const { id } = useParams();
  const [offer, setOffer] = useState(null);
  const [error, setError] = useState(null);

  useEffect(() => {
    axios.get(`/api/offers/${id}`)
      .then(response => {
        setOffer(response.data.data);
      })
      .catch(error => {
        setError(error.message);
      });
  }, [id]);

  if (error) {
    return <div>Error: {error}</div>;
  }

  if (!offer) {
    return <div>Loading...</div>;
  }

  return (
    <div>
      <h1>Offer for Product {offer.product_id}</h1>
      <p>Price: ${offer.price}</p>
      <p>Status: {offer.status}</p>
      <p>Expires At: {new Date(offer.expires_at).toLocaleString()}</p>
    </div>
  );
};

export default OfferDetail;