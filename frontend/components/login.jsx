/**
 * @copyright Popcorn, 2018
 * @author Carmen To
 */
import React from 'react';
import PropTypes from 'prop-types';
// import uuid from 'uuid/v1';
import FlatButton from 'material-ui/FlatButton';
import TextField from 'material-ui/TextField';
import Divider from 'material-ui/Divider';

const fatDividerStyle = { width: '85%', height: '2px', marginTop: '1rem', marginBottom: '1rem' };

class Login extends React.Component {
  state = { username: '', password: '' };

  static propTypes = {
    handleDialogClose: PropTypes.func.isRequired,
    // dispatchLogin: PropTypes.func.isRequired,
    // sessionErrors: PropTypes.array.isRequired,
    // clearErrors: PropTypes.func.isRequired,
    switchDialog: PropTypes.func.isRequired,
  };

  handleFormSubmission = (e) => {
    e.preventDefault();
    // this.props.dispatchLogin(this.state);
  };

  update(field) {
    return (e) => this.setState({ [field]: e.currentTarget.value });
  }

  // get renderErrors() {
  //   if (this.props.sessionErrors === []) {
  //     return;
  //   }
  //   return (
  //     <ul className="session-errors">
  //       {this.props.sessionErrors.map((error) => (
  //         <li key={uuid()} >
  //           {error}
  //         </li>
  //       ))}
  //     </ul>
  //   );
  // }

  render() {
    return (
      <div className="login">
        <form className="login-form" onSubmit={this.handleFormSubmission}>
          <TextField
            fullWidth={true}
            value={this.state.username}
            onChange={this.update('username')}
            hintText="Enter your username"
            floatingLabelText="Username" /><br />
          <TextField
            fullWidth={true}
            value={this.state.password}
            onChange={this.update('password')}
            hintText="Enter your password"
            floatingLabelText="Password"
            type="password" /><br />
          <div className="button-container">
            <FlatButton
              type="login"
              label="Login"
              primary={true}
              keyboardFocused={true} />
            <FlatButton label="Cancel"
              primary={true}
              onClick={this.props.handleDialogClose} />
          </div>
        </form>
        <div className="divider">
          <Divider style={fatDividerStyle} />
        </div>
        <div className="switch-dialog">
          <p>{`Don't have an account?`}
            <FlatButton
              label="Sign up"
              secondary={true}
              onClick={this.props.switchDialog} /></p>
        </div>
      </div>
    );
  }
}

export default Login;
