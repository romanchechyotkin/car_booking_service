import "./rate.css"

export const Rate = (props) => {
    const {rating, comment, user} = props.rate

    return (
        <div className={"rate"} key={rating+user}>
            <div>{user}</div>
            <div>{comment}</div>
            <div>{rating}</div>
        </div>
    )
}