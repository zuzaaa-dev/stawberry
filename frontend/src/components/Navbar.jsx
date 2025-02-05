import React from "react";
import { Link } from "react-router-dom";

function Navbar() {
    return (
        <nav style={{ display: "flex", gap: "1rem", padding: "1rem", backgroundColor: "#f8f9fa" }}>
            <Link to="/" style={{ textDecoration: "none", color: "black" }}>Dashboard</Link>
            <Link to="/login" style={{ textDecoration: "none", color: "black" }}>Login</Link>
            <Link to="/health" style={{ textDecoration: "none", color: "black" }}>Health Status</Link>
            <Link to="/products" style={{ textDecoration: "none", color: "black" }}>Products</Link>
            <Link to="/offers" style={{ textDecoration: "none", color: "black" }}>Offers</Link>
        </nav>
    );
}

export default Navbar;