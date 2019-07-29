
// -*- JavaScript -*-


class UserItem extends React.Component {
  render() {
    return (
      <tr>
        <td scope="row"> {this.props.id}    </td>
        <td scope="row"> {this.props.username} </td>
        <td scope="row"> {this.props.password}  </td>
      </tr>
    );
  }
}

class UsersList extends React.Component {
  constructor(props) {
    super(props);
    this.state = { users: [] };
  }

  componentDidMount() {
    this.serverRequest =
      axios
        .get("/users")
        .then((result) => {
           this.setState({ users: result.data });
        });
  }

  render() {
    const users = this.state.users.map((user, i) => {
      return (
        <UserItem key={i} id={user.ID} username={user.Username} password={user.Password} />
      );
    });

    return (
      <div>
        <table className="table">
        <thead>
          <tr>
            <th scope="col">№</th>
            <th scope="col">username</th>
            <th scope="col">password</th>
          </tr>
        </thead>
        <tbody>
          {users}
        </tbody>
        </table>
      </div>
    );
  }
}

ReactDOM.render( <UsersList/>, document.querySelector("#root"));
