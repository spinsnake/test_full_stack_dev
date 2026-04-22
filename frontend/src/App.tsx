import { Navigate, Route, Routes } from 'react-router-dom';
import GalleryPage from './pages/GalleryPage';
import ManagePage from './pages/ManagePage';

function App() {
  return (
    <Routes>
      <Route path="/" element={<GalleryPage />} />
      <Route path="/manage" element={<ManagePage />} />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}

export default App;
