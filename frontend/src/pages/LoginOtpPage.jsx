import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthProvider";

export default function VerifyOtpPage() {
  const [email, setEmail] = useState("");
  const [otp, setOtp] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const navigate = useNavigate();
  // const { login } = useAuth();

  const onSubmit = async (e) => {
    e.preventDefault();
    const jwtToken = localStorage.getItem("jwtToken");
    const region = localStorage.getItem("region");
    const role = localStorage.getItem("role");

    if (!jwtToken || !region || !role) {
      alert("Missing authentication data. Please log in again.");
      navigate("/login");
      return;
    }

    setIsSubmitting(true);

    try {
      const response = await fetch("", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: jwtToken,
          Region: region,
        },
        body: JSON.stringify({ email, otp }),
      });

      if (response.ok) {
        login(jwtToken, { email }, region);
        navigate("/dashboard");
      } else {
        const errorData = await response.json();
        alert(errorData.error || "OTP verification failed");
      }
    } catch (error) {
      console.error("Verification error:", error);
      alert("An error occurred during verification");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div style={{ display: "flex", alignItems: "center", justifyContent: "center", minHeight: "100vh", backgroundColor: "#f3f4f6" }}>
      <div style={{ backgroundColor: "#fff", padding: "30px", borderRadius: "10px", boxShadow: "0 0 15px rgba(0,0,0,0.1)", maxWidth: "400px", width: "100%" }}>
        <h2 style={{ fontSize: "24px", fontWeight: "600", textAlign: "center", color: "#333", marginBottom: "20px" }}>Verify OTP</h2>
        <form onSubmit={onSubmit} style={{ display: "flex", flexDirection: "column", gap: "15px" }}>
          <div>
            <label style={{ display: "block", fontSize: "14px", fontWeight: "500", marginBottom: "5px", color: "#555" }}>Email</label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              style={{ width: "100%", padding: "10px", border: "1px solid #ccc", borderRadius: "8px", outline: "none", fontSize: "14px" }}
            />
          </div>
          <div>
            <label style={{ display: "block", fontSize: "14px", fontWeight: "500", marginBottom: "5px", color: "#555" }}>OTP Code</label>
            <input
              type="text"
              value={otp}
              onChange={(e) => setOtp(e.target.value)}
              required
              style={{ width: "100%", padding: "10px", border: "1px solid #ccc", borderRadius: "8px", outline: "none", fontSize: "14px" }}
            />
          </div>
          <button
            type="submit"
            disabled={isSubmitting}
            style={{
              width: "100%",
              backgroundColor: isSubmitting ? "#93c5fd" : "#2563eb",
              color: "#fff",
              padding: "10px",
              borderRadius: "8px",
              fontWeight: "600",
              border: "none",
              cursor: isSubmitting ? "not-allowed" : "pointer",
              transition: "background-color 0.3s ease",
            }}
          >
            {isSubmitting ? "Verifying..." : "Verify OTP"}
          </button>
        </form>
        {/* <p style={{ marginTop: "15px", textAlign: "center", fontSize: "13px", color: "#666" }}>
          Didn't receive an OTP? <button style={{ color: "#2563eb", fontWeight: "500", background: "none", border: "none", cursor: "pointer" }} onClick={() => alert("Resend OTP functionality coming soon!")}>Resend OTP</button>
        </p> */}
      </div>
    </div>
  );
}
