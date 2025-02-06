import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import Dashboard from "./pages/Dashboard";
import Login from "./pages/Login";
import HealthStatus from "./pages/HealthStatus";
import ProductList from "./components/ProductList";
import ProductDetail from "./components/ProductDetail";
import OfferList from "./components/OfferList";
import OfferDetail from "./components/OfferDetail";
import { CContainer } from "@coreui/react";

function App() {
  return (
    <Router>
      <CContainer fluid>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/login" element={<Login />} />
          <Route path="/health" element={<HealthStatus />} />
          <Route path="/products" element={<ProductList />} />
          <Route path="/products/:id" element={<ProductDetail />} />
          <Route path="/offers" element={<OfferList />} />
          <Route path="/offers/:id" element={<OfferDetail />} />
        </Routes>
      </CContainer>
    </Router>
  );
}

export default App;