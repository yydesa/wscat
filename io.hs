import Control.Concurrent (forkIO)
import Control.Monad (unless,forever)
import Control.Monad.Trans (liftIO)
import Network.Socket (withSocketsDo)
import Data.ByteString as T
import System.IO(stdin,stdout,stderr)
import Network.WebSockets as WS
import System.IO (hFlush)

app ::WS.ClientApp ()
app conn = do
    -- hPutStrLn stderr "Connected"

    _ <- forkIO $ forever $ do
        msg <- WS.receiveData conn
	-- hPutStrLn stderr $ show msg
	liftIO $ T.hPut stdout msg
	liftIO $ hFlush stdout

    let loop = do
         d <- T.hGetSome stdin 1024
         unless (T.null d) $ WS.sendBinaryData conn d >> loop

    loop
    WS.sendClose conn T.empty

main :: IO ()
main = withSocketsDo $ WS.runClient "127.0.0.1" 8080 "/webssh" app
-- main = withSocketsDo $ WS.runClient "127.0.0.1" 8080 "/echo" app
