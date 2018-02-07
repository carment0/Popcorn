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

class Signup extends React.Component {
  state = { username: '', password: '' };

  static propTypes = {
    handleDialogClose: PropTypes.func.isRequired,
    // dispatchSignup: PropTypes.func.isRequired,
    // sessionErrors: PropTypes.array.isRequired,
    // clearErrors: PropTypes.func.isRequired,
    switchDialog: PropTypes.func.isRequired,
  };

  // componentDidMount() {
  //   this.props.clearErrors();
  // }

  handleFormSubmission = (e) => {
    e.preventDefault();
    // this.props.dispatchSignup(this.state);
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
      <div className="signup">
        <p>By signing up with Popcorn, your personal movie recommendations
           will be saved and updated as you rate more movies.</p>
        <form className="signup-form" onSubmit={this.handleFormSubmission}>
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
          <p>If you do not wish to save your recommendations, click Cancel.</p>
          <div className="button-container">
            <FlatButton
              type="signup"
              label="Signup"
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
          <p>{'Already have a Popcorn account?'}
            <FlatButton
              label="Login"
              secondary={true}
              onClick={this.props.switchDialog} /></p>
        </div>
      </div>
    );
  }
}

export default Signup;
