/**
 * @copyright Popcorn, 2018
 * @author Calvin Feng, Carmen To
 */

import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

import { Nav } from 'react-bootstrap';
import { NavItem } from 'react-bootstrap';
import { Navbar } from 'react-bootstrap';
import Dialog from 'material-ui/Dialog';

import Indicator from '../components/indicator';
import './navigation.scss';
import Login from '../components/login';
import Signup from '../components/signup';

const FormType = {
  SIGN_UP: 'SIGN_UP',
  LOG_IN: 'LOG_IN'
};

const dialogContentStyle = {
  width: '40%',
  minWidth: '350px',
  maxWidth: '500px'
};

class Navigation extends React.Component {
  state = {
    activeKey: 1,
    dialogOpen: false,
    formType: '',
  };

  static propTypes = {
    history: PropTypes.object.isRequired,
    userPreference: PropTypes.object.isRequired
  };

  handleTabSelect = (key) => {
    this.setState({ activeKey: key });

    switch (key) {
      case 1:
        this.setState({ dialogOpen: true, formType: 'SIGN_UP' });
        break;
      case 2:
        window.open('https://github.com/carment0/popcorn');
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

  get dialogTitle() {
    if (this.state.formType === FormType.SIGN_UP) {
      return 'Signup with Popcorn!';
    }
    return 'Login to view your saved recommendations';
  }

  get form() {
    if (this.state.formType === FormType.SIGN_UP) {
      return (
        <Signup
          dispatchSignup={this.props.dispatchSignup}
          handleDialogClose={this.handleDialogClose}
          sessionErrors={this.props.errors}
          clearErrors={this.props.dispatchClearSessionErrors}
          switchDialog={this.createDialogOpenHandler(FormType.LOG_IN)} />
      );
    }
    return (
      <Login
        dispatchLogin={this.props.dispatchLogin}
        handleDialogClose={this.handleDialogClose}
        sessionErrors={this.props.errors}
        clearErrors={this.props.dispatchClearSessionErrors}
        switchDialog={this.createDialogOpenHandler(FormType.SIGN_UP)} />
    );
  }

  render() {
    return (
      <section className="navigation">
        <Navbar inverse collapseOnSelect>
          <Navbar.Header>
            <Navbar.Brand>
              <a href="/#">Popcorn</a>
            </Navbar.Brand>
            <Navbar.Toggle />
          </Navbar.Header>
          <Navbar.Collapse>
            <Nav pullRight onSelect={this.handleTabSelect}>
              <NavItem eventKey={0} id="indicator">
                <Indicator isFetching={this.props.userPreference.isFetching} />
              </NavItem>
              <NavItem eventKey={1}>Recommendations</NavItem>
              <NavItem eventKey={2}>Github</NavItem>
            </Nav>
          </Navbar.Collapse>
        </Navbar>
        <Dialog
          title={this.dialogTitle}
          open={this.state.dialogOpen}
          contentStyle={dialogContentStyle}
          onRequestClose={this.handleDialogClose}>
          {this.form}
        </Dialog>
      </section>
    );
  }
}

const mapReduxStateToProps = (state) => ({ userPreference: state.userPreference });
const mapDispatchToProps = () => ({});

export default connect(mapReduxStateToProps, mapDispatchToProps)(Navigation);
