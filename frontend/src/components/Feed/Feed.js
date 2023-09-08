import {useSelector} from "react-redux";
import {Link} from "react-router-dom";

export const Feed = () => {
    const isVerified = useSelector((state) => state.user.isVerified)
    // const dispatch = useDispatch();

    return (
        <>
            <h1>feed</h1>
            {!isVerified &&
                <div>
                    <h2>you should verify yourself</h2>
                    <Link to={"/verify"}>verify</Link>
                </div>
            }
        </>
    )
}