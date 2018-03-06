/**
 * @copyright Popcorn, 2018
 * @author Calvin Feng, Carmen To
 */

// Thirdparty imports
import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { Nav } from 'react-bootstrap';
import { NavItem } from 'react-bootstrap';
import { Navbar } from 'react-bootstrap';
import Dialog from 'material-ui/Dialog';

// Component imports
import Indicator from '../components/navigation/indicator';
import Login from '../components/navigation/login';
import Signup from '../components/navigation/signup';
import { login, signup, logout, clearSessionErrors, tokenAuthenticate } from '../store/users/session.action';
import { movieRatingsFetch } from '../store/movies/rating.action';

// Style
import './navigation.scss';

const FormType = {
  SIGN_UP: 'SIGN_UP',
  LOG_IN: 'LOG_IN'
};

const dialogContentStyle = {
  width: '40%',
  minWidth: '350px',
  maxWidth: '500px'
};

function capitalizeName(username) {
  return username.charAt(0).toUpperCase() + username.slice(1);
}

class Navigation extends React.Component {
  state = {
    activeKey: 1,
    dialogOpen: false,
    formType: '',
  };

  static propTypes = {
    history: PropTypes.object.isRequired,
    errors: PropTypes.array.isRequired,
    session: PropTypes.object.isRequired,
    dispatchLogin: PropTypes.func.isRequired,
    dispatchSignup: PropTypes.func.isRequired,
    dispatchLogout: PropTypes.func.isRequired,
    dispatchClearSessionErrors: PropTypes.func.isRequired,
    dispatchTokenAuthenticate: PropTypes.func.isRequired,
    dispatchMovieRatingsFetch: PropTypes.func.isRequired
  };

  handleTabSelect = (key) => {
    this.setState({ activeKey: key });
    switch (key) {
      case 1:
        if (this.props.session.currentUser === null) {
          this.setState({ dialogOpen: true, formType: FormType.SIGN_UP });
        } else if (this.props.history.location.pathname !== '/recommend') {
          this.props.history.push('recommend');
        }
        break;
      case 2:
        window.open('https://github.com/carment0/popcorn');
        break;
      case 3:
        this.props.dispatchLogout();
        break;
      case 4:
        this.setState({ dialogOpen: true, formType: FormType.LOG_IN });
        break;
    }
  };

  handleDialogClose = () => {
    this.setState({ dialogOpen: false, formType: '' });
    if (this.props.history.location.pathname !== '/recommend') {
      this.props.history.push('recommend');
    }
  };

  createDialogOpenHandler = (formType) => () => {
    this.setState({ dialogOpen: true, formType: formType });
  };

  /**
   * @return {string}
   */
  get dialogTitle() {
    if (this.state.formType === FormType.SIGN_UP) {
      return 'Signup with Popcorn!';
    }
    return 'Login to view your saved recommendations';
  }

  /**
   * @return {React.Element}
   */
  get form() {
    if (this.state.formType === FormType.SIGN_UP) {
      return (
        <Signup
          dispatchSignup={this.props.dispatchSignup}
          handleDialogClose={this.handleDialogClose}
          sessionErrors={this.props.errors}
          clearErrors={this.props.dispatchClearSessionErrors}
          switchDialog={this.createDialogOpenHandler(FormType.LOG_IN)}
          session={this.props.session} />
      );
    }
    return (
      <Login
        dispatchLogin={this.props.dispatchLogin}
        handleDialogClose={this.handleDialogClose}
        sessionErrors={this.props.errors}
        clearErrors={this.props.dispatchClearSessionErrors}
        switchDialog={this.createDialogOpenHandler(FormType.SIGN_UP)}
        session={this.props.session} />
    );
  }

  /**
   * @return {React.Element}
   */
  get logoutNav() {
    if (this.props.session.currentUser) {
      return (
        <NavItem eventKey={3}>Logout</NavItem>
      );
    }
    return (
      <NavItem eventKey={4}>Login</NavItem>
    );
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.session.currentUser === null) {
      return;
    }

    if (nextProps.session.currentUser !== this.props.session.currentUser) {
      this.props.dispatchMovieRatingsFetch(nextProps.session.currentUser.id);
    }
  }

  componentDidMount() {
    this.props.dispatchTokenAuthenticate();
  }

  render() {
    return (
      <section className="navigation">
        <Navbar inverse collapseOnSelect>
          <Navbar.Header>
            <Navbar.Brand>
              <a href="/#">
                <img
                  className="popcorn-logo"
                  src="images/popcorn-color.svg"
                  alt="logo" />
                <p>
                  {this.props.session.currentUser
                    ? `${capitalizeName(this.props.session.currentUser.username)}'s Popcorn` : 'Popcorn'}
                </p>
              </a>
            </Navbar.Brand>
            <Navbar.Toggle />
          </Navbar.Header>
          <Navbar.Collapse>
            <Nav pullRight onSelect={this.handleTabSelect}>
              <NavItem eventKey={0} id="indicator">
                <Indicator isFetching={false} />
              </NavItem>
              <NavItem eventKey={1}>Recommendations</NavItem>
              <NavItem eventKey={2}>Github</NavItem>
              {this.logoutNav}
            </Nav>
          </Navbar.Collapse>
        </Navbar>
        <Dialog
          title={this.dialogTitle}
          open={this.state.dialogOpen}
          modal={false}
          contentStyle={dialogContentStyle}
          onRequestClose={this.handleDialogClose}>
          {this.form}
        </Dialog>
      </section>
    );
  }
}

const mapReduxStateToProps = (state) => ({
  errors: state.errors.session,
  session: state.session
});

const mapDispatchToProps = (dispatch) => ({
  dispatchLogin: (user) => dispatch(login(user)),
  dispatchSignup: (user) => dispatch(signup(user)),
  dispatchLogout: () => dispatch(logout()),
  dispatchClearSessionErrors: () => dispatch(clearSessionErrors()),
  dispatchTokenAuthenticate: () => dispatch(tokenAuthenticate()),
  dispatchMovieRatingsFetch: (userId) => dispatch(movieRatingsFetch(userId))
});

export default connect(mapReduxStateToProps, mapDispatchToProps)(Navigation);
