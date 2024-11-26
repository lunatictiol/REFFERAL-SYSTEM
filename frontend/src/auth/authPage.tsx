import { useState } from "react";
import { CSSTransition } from "react-transition-group";
import "./styles.css"; // For custom transition classes

interface FormData {
  name?: string;
  email: string;
  password: string;
  referal?: string;
}

const AuthPage: React.FC = () => {
  const [isRegister, setIsRegister] = useState(false);
  const [haveReferal, setHaveReferal] = useState(false);
  const [formData, setFormData] = useState<FormData>({
    name: "",
    email: "",
    password: "",
    referal: "",
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const endpoint = isRegister
      ? "http://localhost:8080/register"
      : "http://localhost:8080/login";

    try {
      const response = await fetch(endpoint, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(formData),
      });

      const result = await response.json();
      console.log(result); // Handle the response
      alert(isRegister ? "Registered Successfully!" : "Logged in Successfully!");
    } catch (error) {
      console.error("Error:", error);
      alert("An error occurred. Please try again.");
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <CSSTransition in appear timeout={300} classNames="fade">
        <div className="w-full max-w-md p-8 space-y-6 bg-white rounded shadow-lg">
          <h2 className="text-2xl font-bold text-center">
            {isRegister ? "Register" : "Login"}
          </h2>

          <form onSubmit={handleSubmit} className="space-y-4">
            {/* Name field (visible only during registration) */}
            <CSSTransition
              in={isRegister}
              timeout={300}
              classNames="slide"
              unmountOnExit
            >
              <div>
                <label
                  htmlFor="name"
                  className="block text-sm font-medium text-gray-700"
                >
                  Name
                </label>
                <input
                  type="text"
                  id="name"
                  name="name"
                  onChange={handleChange}
                  value={formData.name}
                  className="w-full px-3 py-2 mt-1 text-gray-900 border rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
                />
              </div>
            </CSSTransition>

            {/* Referral code toggle button */}
            {isRegister && (
              <div>
                <button
                  type="button"
                  className="text-indigo-600 hover:underline text-sm"
                  onClick={() => setHaveReferal((prev) => !prev)}
                >
                  {haveReferal
                    ? "I don’t have a referral code"
                    : "I have a referral code"}
                </button>
              </div>
            )}

            {/* Referral code field */}
            <CSSTransition
              in={isRegister && haveReferal}
              timeout={300}
              classNames="slide"
              unmountOnExit
            >
              <div>
                <label
                  htmlFor="referal"
                  className="block text-sm font-medium text-gray-700"
                >
                  Referral Code
                </label>
                <input
                  type="text"
                  id="referal"
                  name="referal"
                  onChange={handleChange}
                  value={formData.referal}
                  className="w-full px-3 py-2 mt-1 text-gray-900 border rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
                />
              </div>
            </CSSTransition>

            {/* Email field */}
            <div>
              <label
                htmlFor="email"
                className="block text-sm font-medium text-gray-700"
              >
                Email
              </label>
              <input
                type="email"
                id="email"
                name="email"
                onChange={handleChange}
                value={formData.email}
                required
                className="w-full px-3 py-2 mt-1 text-gray-900 border rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>

            {/* Password field */}
            <div>
              <label
                htmlFor="password"
                className="block text-sm font-medium text-gray-700"
              >
                Password
              </label>
              <input
                type="password"
                id="password"
                name="password"
                onChange={handleChange}
                value={formData.password}
                required
                className="w-full px-3 py-2 mt-1 text-gray-900 border rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>

            {/* Submit button */}
            <button
              type="submit"
              className="w-full px-4 py-2 text-white bg-indigo-600 rounded-md hover:bg-indigo-700 focus:ring-2 focus:ring-indigo-500 transition duration-200"
            >
              {isRegister ? "Register" : "Login"}
            </button>
          </form>

          {/* Toggle between login and register */}
          <p className="text-sm text-center text-gray-600">
            {isRegister ? "Already have an account?" : "Don't have an account?"}{" "}
            <button
              onClick={() => setIsRegister((prev) => !prev)}
              className="font-medium text-indigo-600 hover:underline"
            >
              {isRegister ? "Login" : "Register"}
            </button>
          </p>
        </div>
      </CSSTransition>
    </div>
  );
};

export default AuthPage;
