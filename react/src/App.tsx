import { Routes, Route } from 'react-router-dom'
import { useAuth0 } from "@auth0/auth0-react";
import { PageLoader, LandingPage, HomePage, AboutPage, CallbackPage, ProfilePage } from './pages'
import { AuthenticationGuard, MenuBar } from './components';
import './App.css'

function App() {
  const { isLoading } = useAuth0();

  if (isLoading) {
    return (
      <div className="page-layout">
        <PageLoader />
      </div>
    );
  }

  return (
    <>
      <MenuBar />
      <Routes>
        <Route path="/" element={<LandingPage />} />
        <Route path="/home" element={<AuthenticationGuard component={HomePage} />} />
        <Route path="/profile" element={<AuthenticationGuard component={ProfilePage} />} />
        <Route path="/about" element={<AuthenticationGuard component={AboutPage} />} />
        <Route path="/callback" element={<CallbackPage />} />
        <Route path="*" element={<div>Page Not Found</div>} />
      </Routes>
    </>
  )
}

export default App
