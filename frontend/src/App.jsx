import { Routes, Route, Link, Navigate, useParams } from 'react-router-dom';
import { useEffect, useState } from 'react';
import axios from 'axios';

const API = axios.create({ baseURL: '/api' });

function useAuth() {
  const [token, setToken] = useState(() => localStorage.getItem('token'));
  const login = (t) => { localStorage.setItem('token', t); setToken(t); };
  const logout = () => { localStorage.removeItem('token'); setToken(null); };
  return { token, login, logout };
}

export default function App() {
  const auth = useAuth();
  API.interceptors.request.use(config => {
    if (auth.token) config.headers.Authorization = `Bearer ${auth.token}`;
    return config;
  });

  return (
    <div style={{ padding: 20 }}>
      <nav style={{ marginBottom: 20 }}>
        <Link to="/">Films</Link> | {auth.token ? (
          <>
            <Link to="/add">Add</Link> | <button onClick={auth.logout}>Logout</button>
          </>
        ) : (
          <>
            <Link to="/login">Login</Link> | <Link to="/register">Register</Link>
          </>
        )}
      </nav>
      <Routes>
        <Route path="/" element={<FilmList />} />
        <Route path="/film/:id" element={<FilmDetails />} />
        <Route path="/add" element={auth.token ? <AddFilm /> : <Navigate to="/login" />} />
        <Route path="/login" element={<Login auth={auth} />} />
        <Route path="/register" element={<Register />} />
      </Routes>
    </div>
  );
}

function FilmList() {
  const [films, setFilms] = useState([]);
  useEffect(() => { API.get('/films').then(r => setFilms(r.data)); }, []);
  return (
    <div>
      <h2>Films</h2>
      <ul>
        {films.map(f => (<li key={f.id}><Link to={`/film/${f.id}`}>{f.title}</Link></li>))}
      </ul>
    </div>
  );
}

function FilmDetails() {
  const { id } = useParams();
  const [film, setFilm] = useState(null);
  useEffect(() => { API.get(`/films/${id}`).then(r => setFilm(r.data)); }, [id]);
  if (!film) return <p>Loading...</p>;
  return (
    <div>
      <h2>{film.title}</h2>
      <p>{film.description}</p>
      <p>Rating: {film.rating}</p>
    </div>
  );
}

function AddFilm() {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [release, setRelease] = useState('');
  const submit = () => {
    API.post('/films', { title, description, release_date: new Date(release).toISOString() })
      .then(() => alert('Created'));
  };
  return (
    <div>
      <h2>Add Film</h2>
      <input placeholder="Title" value={title} onChange={e => setTitle(e.target.value)} />
      <input placeholder="Description" value={description} onChange={e => setDescription(e.target.value)} />
      <input type="date" value={release} onChange={e => setRelease(e.target.value)} />
      <button onClick={submit}>Create</button>
    </div>
  );
}

function Login({ auth }) {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const submit = () => {
    API.post('/login', { email, password }).then(r => {
      auth.login(r.data.token);
    });
  };
  return (
    <div>
      <h2>Login</h2>
      <input placeholder="Email" value={email} onChange={e => setEmail(e.target.value)} />
      <input type="password" placeholder="Password" value={password} onChange={e => setPassword(e.target.value)} />
      <button onClick={submit}>Login</button>
    </div>
  );
}

function Register() {
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const submit = () => {
    API.post('/register', { username, email, password }).then(() => alert('registered'));
  };
  return (
    <div>
      <h2>Register</h2>
      <input placeholder="Username" value={username} onChange={e => setUsername(e.target.value)} />
      <input placeholder="Email" value={email} onChange={e => setEmail(e.target.value)} />
      <input type="password" placeholder="Password" value={password} onChange={e => setPassword(e.target.value)} />
      <button onClick={submit}>Register</button>
    </div>
  );
} 