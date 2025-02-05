import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import Dashboard from "./pages/Dashboard";
import Login from "./pages/Login";
import { CContainer } from "@coreui/react";

function App() {
  return (
    <Router>
      <CContainer fluid>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/login" element={<Login />} />
        </Routes>
      </CContainer>
    </Router>
  );
}

export default App;