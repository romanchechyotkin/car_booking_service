import { useEffect, useState } from "react";
import "./navbar.css";
import { userActions } from "../../store/loginUserSlice";
import { useDispatch, useSelector } from "react-redux";
import { Link } from "react-router-dom";
import logoImage from "../../img/logo3.png"
import searchIcon from "../../img/navbar/search_icon.png"
import filterIcon from "../../img/navbar/filter_icon.png"
import settingsIcon from "../../img/navbar/settings.png"
import favoritesIcon from "../../img/navbar/favorites.png"

export const Navbar = () => {
    const dispatch = useDispatch();
    const { isAuth, role, isVerified } = useSelector((state) => state.user);

    const [scroll, setScroll] = useState(false);

    const handleScroll = () => {
    if (window.scrollY > 0) {  
      setScroll(true);
    } else {
      setScroll(false);
    }
  };

  useEffect(() => {
    window.addEventListener('scroll', handleScroll);

    return () => {
      window.removeEventListener('scroll', handleScroll); 
    };
  }, []);
    
    const logout = () => {
        dispatch(userActions.logout());
        localStorage.clear();
    };

    const [searchText, setSearchText] = useState("");
   
    const handleSearchChange = (e) => {
        setSearchText(e.target.value);
    };

    const clearSearch = () => {
        setSearchText("");
    };

    return (
        <header className="navbar">
         
            <div className="navbar_logo">
              <img src={logoImage} alt="logo" className="logo" />
              <Link to="/">uCar</Link>
            </div>
    
            <div className="search-bar-container">
              <div className="search-bar">
                <img src={searchIcon} alt="Search Icon" className="search-img" />
                <input type="text" placeholder="Найдите свою машину тут.." className="search-input" />
                <img src={filterIcon} alt="Filter Icon" className="filter-img" />
              </div>
            </div>
    
            <ul className="navbar_links">
              {isAuth && role !== 'ADMIN' && !isVerified && (
                <>
                  <li>
                    <Link to="/feed">
                      <img src={settingsIcon} width={40} />
                    </Link>
                  </li>
                  <li>
                    <Link to="/feed">
                      <img src={favoritesIcon} width={40} />
                    </Link>
                  </li>
                  <li>
                    <Link to="/verify">Verify</Link>
                  </li>
                </>
              )}
    
              {isAuth && role !== 'ADMIN' && isVerified && (
                <>
                  <Link to="/feed">На главную</Link>
                </>
              )}
    
              {isAuth && role === 'ADMIN' && (
                <>
                  <li>
                    <Link to="/admin">Admin Panel</Link>
                  </li>
                </>
              )}
    
              {!isAuth && (
                <>
                  <li>
                    <Link to="/login">Login</Link>
                  </li>
                  <li>
                    <Link to="/registration">Register</Link>
                  </li>
                </>
              )}
    
              {isAuth && (
                <li>
                  <button onClick={logout}>Logout</button>
                </li>
              )}
            </ul>
    
        </header>
      );
};
