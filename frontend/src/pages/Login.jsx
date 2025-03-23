import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { FaEye, FaEyeSlash } from "react-icons/fa";

export default function LoginPage() {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    email: "",
    password: "",
    region: "north",
  });
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(false);

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const response = await fetch("http://localhost:2426/student/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(formData),
      });

      const data = await response.json();
      if (!response.ok) throw new Error(data.message || "Login failed");

      localStorage.setItem("jwtToken", data.token);
      localStorage.setItem("region", formData.region);

      navigate("/verifyotp");
    } catch (error) {
      setError(error.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={styles.container}>
      <div style={styles.cardWrapper}>
        <div style={styles.headerSection}>
          <h2 style={styles.heading}>Welcome Back!</h2>
          <p style={styles.subText}>Login to access your dashboard</p>
        </div>
        {error && <p style={styles.error}>{error}</p>}
        <form onSubmit={handleSubmit} style={styles.form}>
          <div style={styles.inputGroup}>
            <label style={styles.label}>Email</label>
            <input
              type="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              style={styles.input}
              placeholder="example@email.com"
              required
            />
          </div>

          <div style={styles.inputGroup}>
            <label style={styles.label}>Password</label>
            <div style={styles.passwordWrapper}>
              <input
                type={showPassword ? "text" : "password"}
                name="password"
                value={formData.password}
                onChange={handleChange}
                style={styles.input}
                placeholder="Enter your password"
                required
              />
              <button type="button" onClick={() => setShowPassword(!showPassword)} style={styles.toggleButton}>
                {showPassword ? <FaEyeSlash /> : <FaEye />}
              </button>
            </div>
          </div>

          <div style={styles.inputGroup}>
            <label style={styles.label}>Region</label>
            <select name="region" value={formData.region} onChange={handleChange} style={styles.select}>
              {["north", "south", "east", "west"].map((region) => (
                <option key={region} value={region}>
                  {region.charAt(0).toUpperCase() + region.slice(1)}
                </option>
              ))}
            </select>
          </div>

          <button
            type="submit"
            disabled={loading}
            style={{ ...styles.button, ...(loading ? styles.disabledButton : {}) }}
          >
            {loading ? "Logging in..." : "Login"}
          </button>
        </form>

        <p style={styles.registerText}>
          Don't have an account?
          <a href="/SignUp" style={styles.registerLink}> Register here</a>
        </p>
      </div>
    </div>
  );
}

const styles = {
  container: {
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    minHeight: "100vh",
    background: "linear-gradient(135deg, #f2f7ff, #dce3f0)",
    padding: "20px",
  },
  cardWrapper: {
    backgroundColor: "#fff",
    borderRadius: "20px",
    padding: "40px",
    boxShadow: "0 10px 30px rgba(0,0,0,0.2)",
    width: "100%",
    maxWidth: "500px",
  },
  headerSection: {
    textAlign: "center",
    marginBottom: "25px",
  },
  heading: {
    fontSize: "28px",
    fontWeight: 700,
    color: "#1a1a1a",
  },
  subText: {
    fontSize: "14px",
    color: "#777",
  },
  form: {
    width: "100%",
  },
  inputGroup: {
    marginBottom: "18px",
  },
  label: {
    fontSize: "14px",
    marginBottom: "6px",
    color: "#333",
    display: "block",
  },
  input: {
    width: "100%",
    padding: "12px",
    borderRadius: "10px",
    border: "1px solid #ccc",
    fontSize: "14px",
    transition: "border 0.3s ease",
  },
  select: {
    width: "100%",
    padding: "12px",
    borderRadius: "10px",
    border: "1px solid #ccc",
    fontSize: "14px",
  },
  passwordWrapper: {
    position: "relative",
  },
  toggleButton: {
    position: "absolute",
    top: "50%",
    right: "10px",
    transform: "translateY(-50%)",
    background: "transparent",
    border: "none",
    fontSize: "16px",
    cursor: "pointer",
  },
  button: {
    width: "100%",
    backgroundColor: "#0056d2",
    color: "white",
    padding: "14px",
    fontSize: "16px",
    fontWeight: 600,
    border: "none",
    borderRadius: "10px",
    cursor: "pointer",
    marginTop: "10px",
    transition: "background-color 0.3s ease",
  },
  disabledButton: {
    backgroundColor: "#a5b6d2",
    cursor: "not-allowed",
  },
  registerText: {
    marginTop: "20px",
    textAlign: "center",
    fontSize: "14px",
    color: "#555",
  },
  registerLink: {
    color: "#0056d2",
    textDecoration: "none",
    fontWeight: 600,
    marginLeft: "6px",
  },
  error: {
    backgroundColor: "#ffe0e0",
    color: "#d8000c",
    padding: "10px",
    borderRadius: "8px",
    marginBottom: "15px",
    textAlign: "center",
  },
};
